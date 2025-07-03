package callbacks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

// InfrastructureInfoCallback handles incoming infrastructure information messages from Kafka
type FederationInfrastructureInfoCallback struct {
	services *router.Services
}

// NewFederationInfrastructureInfoCallback creates a new FederationInfrastructureInfoCallback instance
func NewFederationInfrastructureInfoCallback(services *router.Services) *FederationInfrastructureInfoCallback {
	return &FederationInfrastructureInfoCallback{
		services: services,
	}
}

// HandleMessage processes incoming cluster information messages
func (cc *FederationInfrastructureInfoCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("FederationInfrastructureInfoCallback.HandleMessage", func() {
		// Get federations that this federator is a member of
		federations, err := cc.services.FederationService.GetOriginFederations()
		if err != nil {
			log.Println("Error getting federations:", err)
			return
		}

		// Unmarshal the message
		var metricsData map[string]any
		if err := json.Unmarshal(message.Value, &metricsData); err != nil {
			log.Println("Error unmarshalling metrics data:", err)
			return
		}

		// remove msg_id from metricsData
		delete(metricsData, "msg_id")

		log.Printf("Metrics data: %v", metricsData)

		// Send metrics to each federation partner
		for _, federation := range federations {
			if err := cc.sendMetricsToPartner(&federation, metricsData); err != nil {
				log.Printf("Error sending metrics to federation %s: %v", federation.PartnerOP.FederationContextId, err)
			}
		}
	})
}

// sendMetricsToPartner sends metrics data to a specific federation partner
func (cc *FederationInfrastructureInfoCallback) sendMetricsToPartner(federation *models.Federation, metricsData map[string]any) error {
	// Use stored access token from federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// Construct partner endpoint URL for PostMetricsController
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/metrics",
		federation.FederationEndpoint,
		federation.PartnerOP.FederationContextId)

	// Create request payload
	payload := metricsData

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling metrics payload: %v", err)
	}

	// Create HTTP request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Set headers and authentication
	headers := map[string]string{"Content-Type": "application/json"}
	auth := services.NewBearerTokenAuth(accessToken)

	// Make HTTP request via HttpClientService
	resp, err := cc.services.HttpClientService.DoRequest(
		ctx, "POST", partnerEndpoint, bytes.NewReader(payloadBytes), headers, auth)
	if err != nil {
		return fmt.Errorf("error making HTTP request to partner: %v", err)
	}
	defer resp.Body.Close()

	// Handle response
	if resp.StatusCode != 200 {
		return fmt.Errorf("partner returned error status %d", resp.StatusCode)
	}

	log.Printf("Successfully sent metrics to federation partner %s", federation.PartnerOP.FederationContextId)
	return nil
}
