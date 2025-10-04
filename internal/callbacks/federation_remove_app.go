package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type FederationRemoveAppCallback struct {
	services *router.Services
}

func NewFederationRemoveAppCallback(services *router.Services) *FederationRemoveAppCallback {
	return &FederationRemoveAppCallback{
		services: services,
	}
}

func (f *FederationRemoveAppCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("FederationRemoveAppCallback.HandleMessage", func() {
		log.Printf("Received remove app message from topic %s, partition %d, offset %d",
			message.Topic, message.Partition, message.Offset)

		// unmarshal the message
		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		log.Printf("Processing remove app request with message ID: %s", msg["msg_id"])

		msgId := msg["msg_id"].(string)
		f.handleRemoveApp(msgId, msg)
	})
}

func (f *FederationRemoveAppCallback) handleRemoveApp(msgId string, msg map[string]interface{}) {
	appId, ok := msg["app_id"].(string)
	if !ok {
		log.Printf("Error: app_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "app_id is required")
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

	// remove the app on the partner
	log.Printf("Removing app on partner for federation: %s, appId: %s", federationContextId, appId)
	err = f.removeAppOnPartner(&federation, appId)
	if err != nil {
		log.Printf("Error removing app on partner: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", "Error removing app on partner")
		return
	}

	// send the response
	response := map[string]string{
		"msg_id":  msgId,
		"status":  "200",
		"app_id":  appId,
		"message": "App removed on partner successfully",
	}

	_, err = f.services.KafkaClientService.Produce("responses", response)
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}

	log.Printf("Successfully removed app on partner for federation: %s, appId: %s", federation.PartnerOP.FederationContextId, appId)
}

func (f *FederationRemoveAppCallback) removeAppOnPartner(federation *models.Federation, appId string) error {
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// construct the partner's app endpoint URL
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/application/onboarding/app/%s", federation.FederationEndpoint, federation.PartnerOP.FederationContextId, appId)
	log.Printf("Sending app to partner endpoint: %s", partnerEndpoint)

	// create the HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	auth := services.NewBearerTokenAuth(accessToken)

	// send the request to the partner
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

	// check the response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to remove app on partner: %v", resp.Status)
	}

	log.Printf("App removed on partner successfully for federation: %s, appId: %s", federation.PartnerOP.FederationContextId, appId)

	return nil
}
