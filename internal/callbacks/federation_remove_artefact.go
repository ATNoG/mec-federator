package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
)

type FederationArtefactRemoveCallback struct {
	services *router.Services
}

func NewFederationArtefactRemoveCallback(services *router.Services) *FederationArtefactRemoveCallback {
	return &FederationArtefactRemoveCallback{
		services: services,
	}
}

// receives info about an artefact to remove from a certain federation
func (f *FederationArtefactRemoveCallback) HandleMessage(message *sarama.ConsumerMessage) {
	log.Printf("Received remove artefact message from topic %s, partition %d, offset %d", 
		message.Topic, message.Partition, message.Offset)
	
	// unmarshal the message
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Processing remove artefact request with message ID: %s", msg["msg_id"])

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

	// check if an artefact with this appPkgId exists in the federation
	log.Printf("Looking for artefact with appPkgId: %s in federation: %s", appPkgId, federationContextId)
	artefact, err := f.services.ArtefactService.GetArtefactByAppPkgId(federationContextId, appPkgId)
	if err != nil {
		log.Printf("No artefact found with appPkgId %s in federation %s: %v", appPkgId, federationContextId, err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Artefact not found")
		return
	}
	log.Printf("Found artefact with ID: %s", artefact.Id)

	// remove the artefact from the partner operator
	log.Printf("Removing artefact %s from partner operator", artefact.Id)
	err = f.removeArtefactFromPartner(&federation, artefact.Id)
	if err != nil {
		log.Printf("Error removing artefact from partner: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to remove artefact from partner: %v", err))
		return
	}
	log.Printf("Successfully removed artefact from partner operator")

	// remove the artefact locally
	log.Printf("Removing artefact %s from local database", artefact.Id)
	err = f.services.ArtefactService.RemoveArtefact(federationContextId, artefact.Id)
	if err != nil {
		log.Printf("Error removing artefact locally: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", "Failed to remove artefact locally")
		return
	}

	log.Printf("Successfully removed artefact %s from partner operator and locally", appPkgId)

	// send response to kafka with artefact ID
	response := map[string]string{
		"msg_id":      msgId,
		"artefact_id": artefact.Id,
		"status":      "200",
		"message":     "Artefact removed successfully",
	}
	
	_, err = f.services.KafkaClientService.Produce("responses", response)
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}

func (f *FederationArtefactRemoveCallback) removeArtefactFromPartner(federation *models.Federation, artefactId string) error {
	// use the stored access token from the federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// construct the partner's artefact endpoint URL for deletion
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/artefact/%s", federation.FederationEndpoint, federation.PartnerOP.FederationContextId, artefactId)
	log.Printf("Sending delete request to partner endpoint: %s", partnerEndpoint)

	// create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	auth := services.NewBearerTokenAuth(accessToken)

	// make the HTTP DELETE request
	resp, err := f.services.HttpClientService.DoRequest(
		ctx,
		http.MethodDelete,
		partnerEndpoint,
		nil,
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
	log.Printf("Received response from partner with status: %d", resp.StatusCode)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("partner returned error status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Partner response: %s", string(respBody))
	return nil
}
