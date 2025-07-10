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
	"github.com/mankings/mec-federator/internal/utils"
)

type FederationMigrateNodeCallback struct {
	services *router.Services
}

func NewFederationMigrateNodeCallback(services *router.Services) *FederationMigrateNodeCallback {
	return &FederationMigrateNodeCallback{services: services}
}

func (f *FederationMigrateNodeCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("FederationMigrateNodeCallback.HandleMessage", func() {
		log.Printf("Received migrate node message from topic %s, partition %d, offset %d",
			message.Topic, message.Partition, message.Offset)

		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		log.Printf("Processing migrate node request with message ID: %s", msg["msg_id"])

		msgId := msg["msg_id"].(string)
		f.handleMigrateNode(msgId, msg)
	})
}

func (f *FederationMigrateNodeCallback) handleMigrateNode(msgId string, msg map[string]interface{}) {
	// Extract required fields from the message
	federationContextId, ok := msg["federation_context_id"].(string)
	if !ok {
		log.Printf("Error: federation_context_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "federation_context_id is required")
		return
	}

	nsId, ok := msg["ns_id"].(string)
	if !ok {
		log.Printf("Error: ns_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "ns_id is required")
		return
	}

	vnfId, ok := msg["vnf_id"].(string)
	if !ok {
		log.Printf("Error: vnf_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "vnf_id is required")
		return
	}

	kduId, ok := msg["kdu_id"].(string)
	if !ok {
		log.Printf("Error: kdu_id not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "kdu_id is required")
		return
	}

	node, ok := msg["node"].(string)
	if !ok {
		log.Printf("Error: node not found or not a string")
		f.services.KafkaClientService.SendResponse(msgId, "400", "node is required")
		return
	}

	log.Printf("Handling migrate node request - Federation: %s, NsId: %s, VnfId: %s, KduId: %s, Node: %s",
		federationContextId, nsId, vnfId, kduId, node)

	// Get the federation from the database
	federation, err := f.services.FederationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "Federation not found")
		return
	}

	// Get the app instance from the database
	appInstance, err := f.services.AppInstanceService.GetAppInstanceFromNsIdAndVnfId(federationContextId, nsId, vnfId)
	if err != nil {
		log.Printf("Error getting app instance: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "404", "App instance not found")
		return
	}

	// Send the migrate node request to the partner
	log.Printf("Sending migrate node request to partner - Federation: %s, AppInstanceId: %s, NsId: %s, VnfId: %s, KduId: %s, Node: %s",
		federationContextId, appInstance.Id, nsId, vnfId, kduId, node)
	err = f.sendMigrateNodeRequestToPartner(&federation, appInstance.Id, nsId, vnfId, kduId, node)
	if err != nil {
		log.Printf("Error sending migrate node request to partner: %v", err)
		f.services.KafkaClientService.SendResponse(msgId, "500", "Failed to send migrate node request to partner")
		return
	}

	log.Printf("Migrate node request sent to partner successfully")
}

func (f *FederationMigrateNodeCallback) sendMigrateNodeRequestToPartner(federation *models.Federation, appInstanceId string, nsId string, vnfId string, kduId string, node string) error {
	// Use the stored access token from the federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// Construct the partner's enable KDU endpoint URL
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/application/lcm/%s/node/migrate",
		federation.FederationEndpoint, federation.PartnerOP.FederationContextId, appInstanceId)

	// Create the enable request
	request := dto.AppInstanceNodeMigrateRequest{
		NsId:  nsId,
		VnfId: vnfId,
		KduId: kduId,
		Node:  node,
	}

	// Marshal the request
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	authStrategy := services.NewBearerTokenAuth(accessToken)

	// Send the request to the partner
	resp, err := f.services.HttpClientService.DoRequest(
		ctx,
		http.MethodPost,
		partnerEndpoint,
		bytes.NewBuffer(payload),
		headers,
		authStrategy)
	if err != nil {
		return fmt.Errorf("failed to send request to partner: %v", err)
	}

	// status to float
	

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to migrate node: %s", resp.Status)
	}

	log.Printf("Migrate node request sent to partner successfully")

	return nil
}
