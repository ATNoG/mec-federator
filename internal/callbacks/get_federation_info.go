package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/atnog/mec-federator/internal/models"
	"github.com/atnog/mec-federator/internal/router"
	"github.com/atnog/mec-federator/internal/services"
	"github.com/atnog/mec-federator/internal/utils"
)

type GetFederationInfoCallback struct {
	services *router.Services
}

func NewGetFederationInfoCallback(services *router.Services) *GetFederationInfoCallback {
	return &GetFederationInfoCallback{
		services: services,
	}
}

func (g *GetFederationInfoCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("GetFederationInfoCallback.HandleMessage", func() {
		log.Printf("Received get federation info message from topic %s, partition %d, offset %d",
			message.Topic, message.Partition, message.Offset)

		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		log.Printf("Processing get federation info request with message ID: %s", msg["msg_id"])

		msgId := msg["msg_id"].(string)
		g.handleGetFederationInfo(msgId, msg)
	})

}

func (g *GetFederationInfoCallback) handleGetFederationInfo(msgId string, msg map[string]interface{}) {
	// Extract and validate required fields from the message
	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		g.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	// check if the federation exists
	federation, err := g.services.FederationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		g.services.KafkaClientService.SendResponse(msgId, "404", "Federation not found")
		return
	}

	accessToken := federation.OriginOP.AccessToken
	authStrat := services.NewBearerTokenAuth(accessToken.AccessToken)
	headers := map[string]string{"Content-Type": "application/json"}

	// make info request
	getFederationInfoUrl := fmt.Sprintf("%s/federation/v1/ewbi/%s/partner", federation.FederationEndpoint, federation.PartnerOP.FederationContextId)
	log.Printf("Sending get federation info request to: %s", getFederationInfoUrl)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultHTTPTimeout)
	defer cancel()
	resp, err := g.services.HttpClientService.DoRequest(
		ctx,
		http.MethodGet,
		getFederationInfoUrl,
		nil,
		headers,
		authStrat)
	if err != nil {
		log.Printf("Error getting federation info: %v", err)
		g.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to get federation info: %v", err))
		return
	}
	defer resp.Body.Close()

	log.Printf("Received get federation info response with status: %d", resp.StatusCode)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Error getting federation info: %v", resp.Status)
		g.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Partner returned error status %d", resp.StatusCode))
		return
	}

	var federationInfo models.FederationMetaInfo
	err = json.NewDecoder(resp.Body).Decode(&federationInfo)
	if err != nil {
		log.Printf("Error decoding federation info: %v", err)
		g.services.KafkaClientService.SendResponse(msgId, "500", "Failed to decode federation info")
		return
	}

	log.Printf("Federation info: %v", federationInfo)

	// send response to kafka
	err = g.services.KafkaClientService.SendResponse(msgId, "200", "Federation info retrieved successfully")
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}
