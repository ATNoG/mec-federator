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
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
)

type FederationKduDisableCallback struct {
	services *router.Services
}

func NewDisableAppInstanceKDUCallback(services *router.Services) *FederationKduDisableCallback {
	return &FederationKduDisableCallback{
		services: services,
	}
}

func (f *FederationKduDisableCallback) HandleMessage(message *sarama.ConsumerMessage) {
	log.Printf("Received disable KDU message from topic %s, partition %d, offset %d",
		message.Topic, message.Partition, message.Offset)

	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Processing disable KDU request with message ID: %s", msg["msg_id"])

	msgId := msg["msg_id"].(string)

	// Extract required fields from the message
	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	mecAppdId, ok := msg["mec_appd_id"].(string)
	if !ok {
		log.Printf("Error: mec_appd_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "mec_appd_id is required")
		return
	}

	nsId, ok := msg["ns_id"].(string)
	if !ok {
		log.Printf("Error: ns_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "ns_id is required")
		return
	}

	kduId, ok := msg["kdu_id"].(string)
	if !ok {
		log.Printf("Error: kdu_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "kdu_id is required")
		return
	}

	log.Printf("Handling KDU disable request - Federation: %s, NsId: %s, KDU: %s",
		federationContextId, nsId, kduId)

	// Get the federation context
	log.Printf("Retrieving federation with context ID: %s", federationContextId)
	federation, err := f.services.FederationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Federation not found")
		return
	}

	// Get the app pkg from the orchestrator
	appPkg, err := f.services.OrchestratorService.GetAppPkgByMecAppdId(mecAppdId)
	if err != nil {
		log.Printf("Error getting app pkg: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "App pkg not found")
		return
	}

	// Get the federated artefact from the database
	artefact, err := f.services.ArtefactService.GetArtefactByAppPkgId(federationContextId, appPkg.Id.Hex())
	if err != nil {
		log.Printf("Error getting artefact: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Artefact not found")
		return
	}

	// Get the app instance from the database
	appInstance, err := f.services.AppInstanceService.GetAppInstanceFromNsId(federationContextId, artefact.Id, nsId)
	if err != nil {
		log.Printf("Error getting app instance: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "App instance not found")
		return
	}

	// Send disable KDU request to partner
	log.Printf("Sending disable KDU request to partner for KDU: %s", kduId)
	err = f.sendDisableKDURequestToPartner(&federation, appInstance.Id, nsId, kduId)
	if err != nil {
		log.Printf("Error sending disable KDU request to partner: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", fmt.Sprintf("Failed to disable KDU: %v", err))
		return
	}

	// Send success response to kafka
	err = f.services.KafkaClientService.SendResponse(msgId, "200", "KDU disabled successfully")
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}

	log.Printf("Successfully disabled KDU %s for app instance %s", kduId, appInstance.Id)
}

func (f *FederationKduDisableCallback) sendDisableKDURequestToPartner(federation *models.Federation, appInstanceId, nsId, kduId string) error {
	// Use the stored access token from the federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// Construct the partner's disable KDU endpoint URL
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/application/lcm/%s/kdu/disable",
		federation.FederationEndpoint, federation.PartnerOP.FederationContextId, appInstanceId)
	log.Printf("Sending disable KDU request to partner endpoint: %s", partnerEndpoint)

	// Create the disable request
	request := dto.DisableAppInstanceKDURequest{
		KduId: kduId,
		NsId:  nsId,
	}

	// Marshal the request
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	auth := services.NewBearerTokenAuth(accessToken)

	// Make the HTTP request
	resp, err := f.services.HttpClientService.DoRequest(
		ctx,
		http.MethodPatch,
		partnerEndpoint,
		bytes.NewBuffer(payload),
		headers,
		auth)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	log.Printf("Received response from partner with status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("partner returned error status %d", resp.StatusCode)
	}

	log.Printf("Successfully sent disable KDU request to partner for KDU %s", kduId)
	return nil
}
