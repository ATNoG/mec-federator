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
	"github.com/mankings/mec-federator/internal/services"
)

// NewFederationCallback handles incoming new federation messages from Kafka
type NewFederationCallback struct {
	authService        *services.AuthService
	httpClientService  *services.HttpClientService
	kafkaClientService *services.KafkaClientService
	federationService  *services.FederationService
}

// NewNewFederationCallback creates a new NewFederationCallback instance
func NewNewFederationCallback(authService *services.AuthService, httpClientService *services.HttpClientService, kafkaClientService *services.KafkaClientService, federationService *services.FederationService) *NewFederationCallback {
	return &NewFederationCallback{
		authService:        authService,
		httpClientService:  httpClientService,
		kafkaClientService: kafkaClientService,
		federationService:  federationService,
	}
}

// HandleMessage processes incoming new federation messages
func (nc *NewFederationCallback) HandleMessage(message *sarama.ConsumerMessage) {
	// unmarshal the message
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling response message: %v", err)
		return
	}

	federationEndpoint := msg["federation_endpoint"].(string)
	authEndpoint := msg["auth_endpoint"].(string)
	clientId := msg["client_id"].(string)
	clientSecret := msg["client_secret"].(string)

	// fetch access token from the provided auth endpoint
	accessToken, err := nc.authService.FetchAccessTokenFromAuthEndpoint(authEndpoint, clientId, clientSecret)
	if err != nil {
		log.Printf("Error fetching access token: %v", err)
		return
	}

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
		return
	}

	createFederationUrl := fmt.Sprintf("%s/federation/v1/ewbi/partner", federationEndpoint)
	authStrat := services.NewBearerTokenAuth(accessToken.AccessToken)
	headers := map[string]string{"Content-Type": "application/json"}
	resp, err := nc.httpClientService.DoRequest(
		context.Background(),
		http.MethodPost,
		createFederationUrl,
		bytes.NewBuffer(payload),
		headers,
		authStrat)
	if err != nil {
		log.Printf("Error sending federation request: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error creating federation: %v", resp.StatusCode)
		return
	}

	var federationResponseData models.FederationResponseData
	err = json.NewDecoder(resp.Body).Decode(&federationResponseData)
	if err != nil {
		log.Printf("Error decoding federation response: %v", err)
		return
	}

	var federation models.Federation
	federation.PartnerOP = federationResponseData
	federation.OriginOP = federationRequestData
	federation.IsEstablished = true
	federation.AuthEndpoint = authEndpoint
	federation.FederationEndpoint = federationEndpoint
	federation.IsOriginOP = true

	_, err = nc.federationService.CreateFederation(federation)
	if err != nil {
		log.Printf("Error creating federation: %v", err)
		return
	}

	log.Printf("Federation established successfully with partner %s", federationResponseData.PartnerOPFederationId)
	// send response to kafka
	_, err = nc.kafkaClientService.Produce("responses", map[string]string{
		"msg_id": msg["msg_id"].(string),
		"status": "200",
	})
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}
