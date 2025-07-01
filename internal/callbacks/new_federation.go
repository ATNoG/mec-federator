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
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
)

// NewFederationCallback handles incoming new federation messages from Kafka
type NewFederationCallback struct {
	services *router.Services
}

// NewNewFederationCallback creates a new NewFederationCallback instance
func NewNewFederationCallback(services *router.Services) *NewFederationCallback {
	return &NewFederationCallback{
		services: services,
	}
}

// HandleMessage processes incoming new federation messages
func (nc *NewFederationCallback) HandleMessage(message *sarama.ConsumerMessage) {
	log.Printf("Received new federation message from topic %s, partition %d, offset %d",
		message.Topic, message.Partition, message.Offset)

	// unmarshal the message
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Processing new federation request with message ID: %s", msg["msg_id"])

	msgId := msg["msg_id"].(string)

	// Extract and validate required fields from the message
	federationEndpoint, ok := msg["federation_endpoint"].(string)
	if !ok {
		log.Printf("Error: federation_endpoint not found or not a string")
		nc.services.KafkaClientService.SendResponse(msgId, "400", "federation_endpoint is required")
		return
	}

	authEndpoint, ok := msg["auth_endpoint"].(string)
	if !ok {
		log.Printf("Error: auth_endpoint not found or not a string")
		nc.services.KafkaClientService.SendResponse(msgId, "400", "auth_endpoint is required")
		return
	}

	clientId, ok := msg["client_id"].(string)
	if !ok {
		log.Printf("Error: client_id not found or not a string")
		nc.services.KafkaClientService.SendResponse(msgId, "400", "client_id is required")
		return
	}

	clientSecret, ok := msg["client_secret"].(string)
	if !ok {
		log.Printf("Error: client_secret not found or not a string")
		nc.services.KafkaClientService.SendResponse(msgId, "400", "client_secret is required")
		return
	}

	// fetch access token from the provided auth endpoint
	log.Printf("Fetching access token from auth endpoint: %s", authEndpoint)
	accessToken, err := nc.services.AuthService.FetchAccessTokenFromAuthEndpoint(authEndpoint, clientId, clientSecret)
	if err != nil {
		log.Printf("Error fetching access token: %v", err)
		nc.services.KafkaClientService.SendResponse(msgId, "401", "Failed to fetch access token")
		return
	}
	log.Printf("Successfully obtained access token")

	// start federation procedure
	var federationRequestData models.FederationRequestData
	federationRequestData.InitialDate = time.Now()
	federationRequestData.OrigOPFederationId = config.AppConfig.OperatorId
	federationRequestData.OrigOPCountryCode = "351"
	federationRequestData.PartnerStatusLink = fmt.Sprintf("%s%s", federationEndpoint, "/federation/v1/ewbi/partner/status")
	federationRequestData.AccessToken = accessToken

	payload, err := json.Marshal(federationRequestData)
	if err != nil {
		log.Printf("Error marshalling federation request data: %v", err)
		nc.services.KafkaClientService.SendResponse(msgId, "500", "Failed to marshal federation request")
		return
	}

	createFederationUrl := fmt.Sprintf("%s/federation/v1/ewbi/partner", federationEndpoint)
	log.Printf("Sending federation request to: %s", createFederationUrl)
	authStrat := services.NewBearerTokenAuth(accessToken.AccessToken)
	headers := map[string]string{"Content-Type": "application/json"}
	resp, err := nc.services.HttpClientService.DoRequest(
		context.Background(),
		http.MethodPost,
		createFederationUrl,
		bytes.NewBuffer(payload),
		headers,
		authStrat)
	if err != nil {
		log.Printf("Error sending federation request: %v", err)
		nc.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to send federation request: %v", err))
		return
	}

	defer resp.Body.Close()

	log.Printf("Received federation response with status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error creating federation: %v", resp.StatusCode)
		nc.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Partner returned error status %d", resp.StatusCode))
		return
	}

	var federationResponseData models.FederationResponseData
	err = json.NewDecoder(resp.Body).Decode(&federationResponseData)
	if err != nil {
		log.Printf("Error decoding federation response: %v", err)
		nc.services.KafkaClientService.SendResponse(msgId, "500", "Failed to decode federation response")
		return
	}

	var federation models.Federation
	federation.PartnerOP = federationResponseData
	federation.OriginOP = federationRequestData
	federation.AuthEndpoint = authEndpoint
	federation.FederationEndpoint = federationEndpoint
	federation.IsEstablished = true
	federation.IsOriginOP = true

	// create federation object locally
	log.Printf("Saving federation locally for partner: %s", federationResponseData.PartnerOPFederationId)
	_, err = nc.services.FederationService.CreateFederation(federation)
	if err != nil {
		log.Printf("Error creating federation: %v", err)
		nc.services.KafkaClientService.SendResponse(msgId, "500", "Failed to save federation locally")
		return
	}

	log.Printf("Federation established successfully with partner %s", federationResponseData.PartnerOPFederationId)

	// send response to kafka
	err = nc.services.KafkaClientService.SendResponse(msgId, "200", "Federation established successfully")
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}
