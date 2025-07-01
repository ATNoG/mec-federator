package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/services"
)

type RemoveFederationCallback struct {
	authService        *services.AuthService
	httpClientService  *services.HttpClientService
	kafkaClientService *services.KafkaClientService
	federationService  *services.FederationService
}

func NewRemoveFederationCallback(authService *services.AuthService, httpClientService *services.HttpClientService, kafkaClientService *services.KafkaClientService, federationService *services.FederationService) *RemoveFederationCallback {
	return &RemoveFederationCallback{
		authService:        authService,
		httpClientService:  httpClientService,
		kafkaClientService: kafkaClientService,
		federationService:  federationService,
	}
}

func (r *RemoveFederationCallback) HandleMessage(message *sarama.ConsumerMessage) {
	// unmarshal the message
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling response message: %v", err)
		return
	}

	federationContextId := msg["federation_context_id"].(string)

	// get the federation from the database
	federation, err := r.federationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		return
	}

	accessToken := federation.OriginOP.AccessToken
	authStrat := services.NewBearerTokenAuth(accessToken.AccessToken)
	headers := map[string]string{"Content-Type": "application/json"}

	// remove the federation
	removeFederationUrl := fmt.Sprintf("%s/federation/v1/ewbi/%s/partner", federation.FederationEndpoint, federation.PartnerOP.FederationContextId)
	resp, err := r.httpClientService.DoRequest(
		context.TODO(),
		http.MethodDelete,
		removeFederationUrl,
		nil,
		headers,
		authStrat)
	if err != nil {
		log.Printf("Error removing federation: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error removing federation: %v", resp.Status)
		return
	}

	// delete the federation from the database
	err = r.federationService.DeleteFederation(federationContextId)
	if err != nil {
		log.Printf("Error deleting federation from database: %v", err)
		return
	}

	log.Printf("Federation removed successfully")
	// send response to kafka
	_, err = r.kafkaClientService.Produce("responses", map[string]string{
		"msg_id": msg["msg_id"].(string),
		"status": "200",
	})
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}
