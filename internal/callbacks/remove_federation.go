package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
)

type RemoveFederationCallback struct {
	services *router.Services
}

func NewRemoveFederationCallback(services *router.Services) *RemoveFederationCallback {
	return &RemoveFederationCallback{
		services: services,
	}
}

func (r *RemoveFederationCallback) HandleMessage(message *sarama.ConsumerMessage) {
	log.Printf("Received remove federation message from topic %s, partition %d, offset %d", 
		message.Topic, message.Partition, message.Offset)
	
	// unmarshal the message
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Processing remove federation request with message ID: %s", msg["msg_id"])

	msgId := msg["msg_id"].(string)

	// Extract and validate required fields from the message
	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		r.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	// get the federation from the database
	log.Printf("Retrieving federation with context ID: %s", federationContextId)
	federation, err := r.services.FederationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		r.services.KafkaClientService.SendResponse(msgId, "404", "Federation not found")
		return
	}
	log.Printf("Successfully retrieved federation for partner: %s", federation.PartnerOP.PartnerOPFederationId)

	accessToken := federation.OriginOP.AccessToken
	authStrat := services.NewBearerTokenAuth(accessToken.AccessToken)
	headers := map[string]string{"Content-Type": "application/json"}

	// remove the federation
	removeFederationUrl := fmt.Sprintf("%s/federation/v1/ewbi/%s/partner", federation.FederationEndpoint, federation.PartnerOP.FederationContextId)
	log.Printf("Sending remove federation request to: %s", removeFederationUrl)
	resp, err := r.services.HttpClientService.DoRequest(
		context.TODO(),
		http.MethodDelete,
		removeFederationUrl,
		nil,
		headers,
		authStrat)
	if err != nil {
		log.Printf("Error removing federation: %v", err)
		r.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to remove federation from partner: %v", err))
		return
	}

	log.Printf("Received remove federation response with status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error removing federation: %v", resp.Status)
		r.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Partner returned error status %d", resp.StatusCode))
		return
	}

	// delete the federation from the database
	log.Printf("Deleting federation from local database: %s", federationContextId)
	err = r.services.FederationService.DeleteFederation(federationContextId)
	if err != nil {
		log.Printf("Error deleting federation from database: %v", err)
		r.services.KafkaClientService.SendResponse(msgId, "500", "Failed to delete federation from local database")
		return
	}

	log.Printf("Federation removed successfully")
	// send response to kafka
	err = r.services.KafkaClientService.SendResponse(msgId, "200", "Federation removed successfully")
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}
