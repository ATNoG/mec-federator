package callbacks

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type GetAppiInfoCallback struct {
	services *router.Services
}

func NewGetAppiInfoCallback(services *router.Services) *GetAppiInfoCallback {
	return &GetAppiInfoCallback{
		services: services,
	}
}

func (g *GetAppiInfoCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("GetAppiInfoCallback.HandleMessage", func() {
		utils.SendResultsMessage(utils.ResultsMessage{
			Name:    "oo-init-get-appi-info",
			Message: "reset",
			Value:   nil,
		})

		log.Printf("Received get app instance info message from topic %s, partition %d, offset %d",
			message.Topic, message.Partition, message.Offset)

		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		log.Printf("Processing get app instance info request with message ID: %s", msg["msg_id"])

		msgId := msg["msg_id"].(string)
		g.handleGetAppiInfo(msgId, msg)

		utils.SendResultsMessage(utils.ResultsMessage{
			Name:    "oo-done-get-appi-info",
			Message: "",
			Value:   nil,
		})
	})

}

func (g *GetAppiInfoCallback) handleGetAppiInfo(msgId string, msg map[string]interface{}) {
	// Extract and validate required fields from the message
	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		g.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	appInstanceId, ok := msg["app_instance_id"].(string)
	if !ok {
		log.Printf("Error: app_instance_id not found or not a string")
		g.services.KafkaClientService.SendResponse(msgId, "400", "app_instance_id is required")
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
	getAppiInfoUrl := fmt.Sprintf("%s/federation/v1/ewbi/%s/application/lcm/%s", federation.FederationEndpoint, federation.PartnerOP.FederationContextId, appInstanceId)
	log.Printf("Sending get app instance info request to: %s", getAppiInfoUrl)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := g.services.HttpClientService.DoRequest(
		ctx,
		http.MethodGet,
		getAppiInfoUrl,
		nil,
		headers,
		authStrat)
	if err != nil {
		log.Printf("Error getting app instance info: %v", err)
		g.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to get app instance info: %v", err))
		return
	}
	defer resp.Body.Close()

	log.Printf("Received get app instance info response with status: %d", resp.StatusCode)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Error getting app instance info: %v", resp.Status)
		g.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Partner returned error status %d", resp.StatusCode))
		return
	}

	var appInstanceInfo models.AppInstance
	err = json.NewDecoder(resp.Body).Decode(&appInstanceInfo)
	if err != nil {
		log.Printf("Error decoding app instance info: %v", err)
		g.services.KafkaClientService.SendResponse(msgId, "500", "Failed to decode app instance info")
		return
	}

	log.Printf("App instance info: %v", appInstanceInfo)

	// send response to kafka
	err = g.services.KafkaClientService.SendResponse(msgId, "200", "App instance info retrieved successfully")
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}
