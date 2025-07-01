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
)

type FederationRemoveAppiCallback struct {
	services *router.Services
}

func NewFederationRemoveAppiCallback(services *router.Services) *FederationRemoveAppiCallback {
	return &FederationRemoveAppiCallback{
		services: services,
	}
}

func (f *FederationRemoveAppiCallback) HandleMessage(message *sarama.ConsumerMessage) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	msgId := msg["msg_id"].(string)

	// Extract and validate required fields from the message
	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	appInstanceId, ok := msg["app_instance_id"].(string)
	if !ok {
		log.Printf("Error: app_instance_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "app_instance_id is required")
		return
	}

	federation, err := f.services.FederationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Federation not found")
		return
	}

	_, err = f.services.AppInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("Error getting app instance: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "App instance not found")
		return
	}

	err = f.sendDeleteRequestToPartner(&federation, appInstanceId)
	if err != nil {
		log.Printf("Error sending delete request to partner: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to delete app instance from partner: %v", err))
		return
	}

	err = f.services.AppInstanceService.RemoveAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("Error deleting app instance from local database: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", "Failed to delete app instance from local database")
		return
	}

	log.Printf("Successfully deleted app instance %s from partner and local database", appInstanceId)

	err = f.services.KafkaClientService.SendResponse(msgId, "200", "App instance removed successfully")
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}

func (f *FederationRemoveAppiCallback) sendDeleteRequestToPartner(federation *models.Federation, appInstanceId string) error {
	accessToken := federation.OriginOP.AccessToken.AccessToken

	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/application/lcm/%s",
		federation.FederationEndpoint,
		federation.PartnerOP.FederationContextId,
		appInstanceId)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	auth := services.NewBearerTokenAuth(accessToken)

	resp, err := f.services.HttpClientService.DoRequest(
		ctx,
		http.MethodDelete,
		partnerEndpoint,
		nil,
		headers,
		auth)
	if err != nil {
		return fmt.Errorf("failed to send HTTP delete request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("partner returned error status %d", resp.StatusCode)
	}

	log.Printf("Partner confirmed deletion of app instance %s", appInstanceId)
	return nil
}
