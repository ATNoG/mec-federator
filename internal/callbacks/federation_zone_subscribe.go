package callbacks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/IBM/sarama"
	"github.com/atnog/mec-federator/internal/models"
	"github.com/atnog/mec-federator/internal/router"
	"github.com/atnog/mec-federator/internal/services"
	"github.com/atnog/mec-federator/internal/utils"
)

type FederationZoneSubscribeCallback struct {
	services *router.Services
}

func NewFederationZoneSubscribeCallback(services *router.Services) *FederationZoneSubscribeCallback {
	return &FederationZoneSubscribeCallback{
		services: services,
	}
}

func (f *FederationZoneSubscribeCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("FederationZoneSubscribeCallback.HandleMessage", func() {
		utils.SendResultsMessage(utils.ResultsMessage{
			Name:    "oo-init-subscribe-zone",
			Message: "reset",
			Value:   nil,
		})

		log.Printf("Received subscribe zone message from topic %s, partition %d, offset %d",
			message.Topic, message.Partition, message.Offset)

		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		log.Printf("Processing subscribe zone request with message ID: %s", msg["msg_id"])

		msgId := msg["msg_id"].(string)
		f.handleSubscribeZone(msgId, msg)

		utils.SendResultsMessage(utils.ResultsMessage{
			Name:    "oo-done-subscribe-zone",
			Message: "",
			Value:   nil,
		})
	})
}

func (f *FederationZoneSubscribeCallback) handleSubscribeZone(msgId string, msg map[string]interface{}) {
	log.Printf("Processing subscribe zone request with message ID: %s", msgId)

	// get the zone id from the message
	zoneId, ok := msg["zone_id"].(string)
	if !ok {
		log.Printf("Error: zone_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "zone_id is required")
		return
	}

	// get the federation context id from the message
	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	// get the federation from the database
	log.Printf("Retrieving federation with context ID: %s", federationContextId)
	federation, err := f.services.FederationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Federation not found")
		return
	}

	// check if zone is in available zones
	availableZones := federation.PartnerOP.OfferedAvailabilityZones

	zoneFound := slices.ContainsFunc(availableZones, func(zone models.ZoneDetails) bool {
		return zone.ZoneId == zoneId
	})

	if !zoneFound {
		log.Printf("Zone is not in partner offered zones: %s", zoneId)
		f.services.KafkaClientService.SendResponse(msgId, "400", "Zone is not in partner's offered availability zones")
		return
	}

	// send subscribe request to partner
	log.Printf("Sending subscribe zone request to partner for zone: %s", zoneId)
	err = f.sendSubscribeRequestToPartner(&federation, zoneId)
	if err != nil {
		log.Printf("Error sending subscribe zone request to partner: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to subscribe to zone: %v", err))
		return
	}

	// send success response to kafka
	err = f.services.KafkaClientService.SendResponse(msgId, "200", "Zone subscribed successfully")
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}

	log.Printf("Successfully subscribed to zone %s for federation %s", zoneId, federationContextId)
}

func (f *FederationZoneSubscribeCallback) sendSubscribeRequestToPartner(federation *models.Federation, zoneId string) error {
	// use the stored access token from the federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// construct the partner's subscribe zone endpoint URL
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/zones",
		federation.FederationEndpoint, federation.PartnerOP.FederationContextId)

	// create the subscribe request
	request := models.ZoneRegistrationRequestData{
		AcceptedAvailabilityZones: []string{zoneId},
		AvailZoneNotifLink:        "link",
	}

	// marshal the request
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), LongHTTPTimeout)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	auth := services.NewBearerTokenAuth(accessToken)

	// make the HTTP request
	resp, err := f.services.HttpClientService.DoRequest(
		ctx,
		http.MethodPost,
		partnerEndpoint,
		bytes.NewBuffer(payload),
		headers,
		auth)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("partner returned error status %d", resp.StatusCode)
	}

	log.Printf("Successfully sent subscribe zone request to partner for zone %s", zoneId)
	return nil
}
