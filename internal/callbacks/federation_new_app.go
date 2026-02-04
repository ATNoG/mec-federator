package callbacks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/atnog/mec-federator/internal/models"
	"github.com/atnog/mec-federator/internal/models/dto"
	"github.com/atnog/mec-federator/internal/router"
	"github.com/atnog/mec-federator/internal/services"
	"github.com/atnog/mec-federator/internal/utils"
	"github.com/google/uuid"
)

type FederationNewAppCallback struct {
	services *router.Services
}

func NewFederationNewAppCallback(services *router.Services) *FederationNewAppCallback {
	return &FederationNewAppCallback{
		services: services,
	}
}

func (f *FederationNewAppCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("FederationNewAppCallback.HandleMessage", func() {
		utils.SendResultsMessage(utils.ResultsMessage{
			Name:    "oo-init-onboard-app",
			Message: "reset",
			Value:   nil,
		})

		log.Printf("Received new app message from topic %s, partition %d, offset %d",
			message.Topic, message.Partition, message.Offset)

		// unmarshal the message
		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		log.Printf("Processing new app request with message ID: %s", msg["msg_id"])

		msgId := msg["msg_id"].(string)
		f.handleNewApp(msgId, msg)

		utils.SendResultsMessage(utils.ResultsMessage{
			Name:    "oo-done-onboard-app",
			Message: "",
			Value:   nil,
		})
	})
}

func (f *FederationNewAppCallback) handleNewApp(msgId string, msg map[string]interface{}) {
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

	// get the local artefact linked to the app package
	log.Printf("Retrieving artefact with ID: %s", appPkgId)
	artefact, err := f.services.ArtefactService.GetArtefactByAppPkgId(federationContextId, appPkgId)
	if err != nil {
		log.Printf("Error getting artefact: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Artefact not found")
		return
	}

	// request app onboard to the partner operator
	log.Printf("Requesting app onboard to partner operator for artefact: %s", artefact.Id)
	appId, err := f.onboardAppOnPartner(&federation, &artefact)
	if err != nil {
		log.Printf("Error onboarding app to partner operator: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", "Failed to onboard app to partner operator")
		return
	}

	// make local application object
	application := models.Application{
		Id:                  appId,
		FederationContextId: federationContextId,
		AppPkgId:            "",
		Name:                artefact.Name + "-fed-app",
		Description:         "not-local",
		ArtefactId:          artefact.Id,
	}

	// save the application to the database
	log.Printf("Saving application to database: %v", application)
	err = f.services.ApplicationService.CreateApplication(application)
	if err != nil {
		log.Printf("Error saving application to database: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", "Failed to save application to database")
		return
	}

	// send response to kafka
	response := map[string]string{
		"msg_id":  msgId,
		"status":  "200",
		"app_id":  appId,
		"message": "App onboarded to partner operator successfully",
	}
	_, err = f.services.KafkaClientService.Produce("responses", response)
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}

	log.Printf("Successfully onboarded app to partner operator for federation: %s, appId: %s", federation.PartnerOP.FederationContextId, appId)
}

func (f *FederationNewAppCallback) onboardAppOnPartner(federation *models.Federation, artefact *models.Artefact) (string, error) {
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// construct the partner's app endpoint URL
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/application/onboarding", federation.FederationEndpoint, federation.PartnerOP.FederationContextId)
	log.Printf("Sending app to partner endpoint: %s", partnerEndpoint)

	// create the request body
	request := dto.OnboardApplicationRequest{
		AppId:      uuid.New().String(),
		AppName:    "fed-app",
		ArtefactId: artefact.Id,
	}

	// marshal the request
	payload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// create the HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), LongHTTPTimeout)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	auth := services.NewBearerTokenAuth(accessToken)

	// send the request to the partner operator
	resp, err := f.services.HttpClientService.DoRequest(
		ctx,
		http.MethodPost,
		partnerEndpoint,
		bytes.NewBuffer(payload),
		headers,
		auth)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// check the response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to onboard app to partner operator: %v", resp.Status)
	}

	log.Printf("Successfully onboarded app to partner operator for federation: %s, appId: %s", federation.PartnerOP.FederationContextId, request.AppId)
	return request.AppId, nil
}
