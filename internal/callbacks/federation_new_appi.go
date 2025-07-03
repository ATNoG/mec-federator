package callbacks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type FederationAppiNewCallback struct {
	services *router.Services
}

func NewFederationAppiNewCallback(services *router.Services) *FederationAppiNewCallback {
	return &FederationAppiNewCallback{
		services: services,
	}
}

func (f *FederationAppiNewCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("FederationAppiNewCallback.HandleMessage", func() {
		log.Printf("Received new app instance message from topic %s, partition %d, offset %d",
			message.Topic, message.Partition, message.Offset)

		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		log.Printf("Processing new app instance request with message ID: %s", msg["msg_id"])

		msgId := msg["msg_id"].(string)
		f.handleNewAppInstance(msgId, msg)
	})
}

func (f *FederationAppiNewCallback) handleNewAppInstance(msgId string, msg map[string]interface{}) {
	// Extract and validate required fields from the message
	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	appPkgId, ok := msg["app_pkg_id"].(string)
	if !ok {
		log.Printf("Error: app_pkg_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "app_pkg_id is required")
		return
	}

	vimId, ok := msg["vim_id"].(string)
	if !ok {
		log.Printf("Error: vim_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "vim_id is required")
		return
	}

	config, ok := msg["config"].(string)
	if !ok {
		log.Printf("Error: config not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "config is required")
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

	// check if the artefact exists in the given federation context
	log.Printf("Looking for artefact with app package ID: %s in federation: %s", appPkgId, federationContextId)
	artefact, err := f.services.ArtefactService.GetArtefactByAppPkgId(federationContextId, appPkgId)
	if err != nil {
		log.Printf("Error getting artefact: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Artefact not found")
		return
	}

	// get the partner zone from the vimId
	var zoneId string
	for _, zone := range federation.PartnerOP.OfferedAvailabilityZones {
		if zone.VimId == vimId {
			zoneId = zone.ZoneId
			break
		}
	}

	if zoneId == "" {
		log.Printf("No zone found for vimId: %s", vimId)
		f.services.KafkaClientService.SendResponse(msgId, "404", fmt.Sprintf("Zone not found for vimId: %s", vimId))
		return
	}

	// make request to send to partner
	var request dto.InstantiateApplicationRequest
	request.TransactionId = uuid.New().String()
	request.AppId = artefact.Id
	request.AppProviderId = artefact.AppProviderId
	request.AppVersion = artefact.VersionInfo
	request.ZoneInfo = zoneId
	request.Config = config
	request.AppInstCallbackLink = "callback.link"

	// send request to partner
	appInstanceId, nsId, err := f.sendAppInstanceRequestToPartner(&federation, &request)
	if err != nil {
		log.Printf("Error sending app instance request to partner: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to create app instance: %v", err))
		return
	}

	// make AppInstance object
	appInstance := models.AppInstance{
		Id:                  appInstanceId,
		FederationContextId: federationContextId,
		ArtefactId:          artefact.Id,
		Name:                "federated-instance",
		Description:         "not-local",
		AppPkgId:            "", // empty if the app instance is running on a partner zone
		NsId:                nsId,
	}

	// save appinstance to database
	err = f.services.AppInstanceService.RegisterAppInstance(appInstance)
	if err != nil {
		log.Printf("Error saving app instance locally: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", "Failed to save app instance locally")
		return
	}

	log.Printf("Successfully created app instance %s with partner operator and saved locally", appInstanceId)

	// send response to kafka with app instance ID
	response := map[string]string{
		"msg_id":          msgId,
		"app_instance_id": appInstance.Id,
		"ns_id":           nsId,
		"status":          "201",
		"message":         "App instance created successfully",
	}

	_, err = f.services.KafkaClientService.Produce("responses", response)
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}

func (f *FederationAppiNewCallback) sendAppInstanceRequestToPartner(federation *models.Federation, request *dto.InstantiateApplicationRequest) (string, string, error) {
	// use the stored access token from the federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// construct the partner's app instance endpoint URL
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/application/lcm", federation.FederationEndpoint, federation.PartnerOP.FederationContextId)

	// marshal the request
	payload, err := json.Marshal(request)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	auth := services.NewBearerTokenAuth(accessToken)

	// make the HTTP request
	resp, err := f.services.HttpClientService.DoRequest(
		ctx,
		http.MethodPost,
		partnerEndpoint,
		bytes.NewBuffer(payload),
		headers,
		auth)
	if err != nil {
		return "", "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("partner returned error status %d", resp.StatusCode)
	}

	// parse response to get app instance ID
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode response: %v", err)
	}

	appInstanceId, ok := response["appInstanceId"].(string)
	if !ok {
		return "", "", fmt.Errorf("appInstanceId not found in partner response")
	}

	nsId, ok := response["nsId"].(string)
	if !ok {
		return "", "", fmt.Errorf("nsId not found in partner response")
	}

	log.Printf("Partner response: app instance created with ID %s", appInstanceId)
	return appInstanceId, nsId, nil
}
