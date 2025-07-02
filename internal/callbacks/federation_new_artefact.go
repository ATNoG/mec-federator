package callbacks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
)

type FederationArtefactNewCallback struct {
	services *router.Services
}

func NewFederationArtefactNewCallback(services *router.Services) *FederationArtefactNewCallback {
	return &FederationArtefactNewCallback{
		services: services,
	}
}

// receives info about an artefact to make available to a certain federation
func (f *FederationArtefactNewCallback) HandleMessage(message *sarama.ConsumerMessage) {
	log.Printf("Received new artefact message from topic %s, partition %d, offset %d",
		message.Topic, message.Partition, message.Offset)

	// unmarshal the message
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Processing new artefact request with message ID: %s", msg["msg_id"])

	msgId := msg["msg_id"].(string)

	// Extract and validate required fields from the message
	appPkgId, ok := msg["app_pkg_id"].(string)
	if !ok {
		log.Printf("Error: app_pkg_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "app_pkg_id is required")
		return
	}

	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	// get the federation context
	log.Printf("Retrieving federation with context ID: %s", federationContextId)
	federation, err := f.services.FederationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Federation not found")
		return
	}

	// check if the artefact already exists
	artefact, err := f.services.ArtefactService.GetArtefactByAppPkgId(federationContextId, appPkgId)
	if err == nil {
		log.Printf("Artefact already exists: %s", artefact.Id)
		f.services.KafkaClientService.SendResponse(msgId, "200", "Artefact already exists")
		return
	}

	// get the app package from the orchestrator
	log.Printf("Retrieving app package from orchestrator: %s", appPkgId)
	appPkg, err := f.services.OrchestratorService.GetAppPkg(appPkgId)
	if err != nil {
		log.Printf("Error getting app package from orchestrator: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "App package not found")
		return
	}
	log.Printf("Successfully retrieved app package: %s", appPkg.Name)

	// create the artefact
	artefact = models.Artefact{
		Id:                  uuid.New().String(),
		FederationContextId: federationContextId,
		AppPkgId:            appPkgId,
		AppProviderId:       appPkg.Provider,
		Name:                appPkg.Name,
		Description:         appPkg.Description,
		VersionInfo:         appPkg.Version,
		VirtType:            models.CONTAINER_TYPE,
		DescriptorType:      models.HELM,
		FileFormat:          models.TARGZ,
		ArtefactFile:        &appPkg.AppD,
	}

	// send the artefact to the partner operator
	log.Printf("Sending artefact %s to partner operator", artefact.Id)
	err = f.sendArtefactToPartner(&federation, &artefact)
	if err != nil {
		log.Printf("Error sending artefact to partner: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to send artefact to partner: %v", err))
		return
	}
	log.Printf("Successfully sent artefact to partner operator")

	// save the artefact locally
	log.Printf("Saving artefact %s to local database", artefact.Id)
	artefact.FederationContextId = federationContextId
	err = f.services.ArtefactService.SaveArtefact(artefact)
	if err != nil {
		log.Printf("Error saving artefact locally: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", "Failed to save artefact locally")
		return
	}

	log.Printf("Successfully sent artefact %s to partner operator and saved locally", appPkgId)

	// send response to kafka with artefact ID
	response := map[string]string{
		"msg_id":      msgId,
		"artefact_id": artefact.Id,
		"status":      "200",
		"message":     "Artefact created and sent successfully",
	}

	_, err = f.services.KafkaClientService.Produce("responses", response)
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}

func (f *FederationArtefactNewCallback) sendArtefactToPartner(federation *models.Federation, artefact *models.Artefact) error {
	// use the stored access token from the federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// construct the partner's artefact endpoint URL
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/artefact", federation.FederationEndpoint, federation.PartnerOP.FederationContextId)
	log.Printf("Sending artefact to partner endpoint: %s", partnerEndpoint)

	// create temporary file for the artefact
	tempFile, err := os.CreateTemp("", "artefact_*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// write binary data to temporary file
	if _, err := tempFile.Write(*artefact.ArtefactFile); err != nil {
		return fmt.Errorf("failed to write artefact data to temporary file: %v", err)
	}
	tempFile.Close()

	// reopen file for reading
	file, err := os.Open(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open temporary file: %v", err)
	}
	defer file.Close()

	// create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// add form fields
	writer.WriteField("artefactId", artefact.Id)
	writer.WriteField("appProviderId", artefact.AppProviderId)
	writer.WriteField("artefactName", artefact.Name)
	writer.WriteField("artefactVersionInfo", artefact.VersionInfo)
	writer.WriteField("artefactDescription", artefact.Description)
	writer.WriteField("artefactVirtType", string(artefact.VirtType))
	writer.WriteField("artefactDescriptorType", string(artefact.DescriptorType))
	writer.WriteField("artefactFileFormat", string(artefact.FileFormat))
	writer.WriteField("artefactFileName", artefact.FileName)
	writer.WriteField("repoType", string(models.UPLOAD))

	// add file attachment
	fileName := fmt.Sprintf("%s.tar.gz", artefact.FileName)
	if artefact.FileName != "" {
		fileName = artefact.FileName
	}

	part, err := writer.CreateFormFile("artefactFile", fileName)
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file to form: %v", err)
	}

	writer.Close()

	// create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	auth := services.NewBearerTokenAuth(accessToken)

	// make the HTTP request
	resp, err := f.services.HttpClientService.DoRequest(
		ctx,
		http.MethodPost,
		partnerEndpoint,
		&buf,
		headers,
		auth)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("partner returned error status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Partner response: %s", string(respBody))
	return nil
}
