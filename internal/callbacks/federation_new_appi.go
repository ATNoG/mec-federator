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
	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/services"
)

type FederationAppiNewCallback struct {
	services *router.Services
}

func NewFederationAppiNewCallback(services *router.Services) *FederationAppiNewCallback {
	return &FederationAppiNewCallback{
		services: services,
	}
}

func (f *FederationAppiNewCallback) HandleMessage(message *sarama.ConsumerMessage) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling response message: %v", err)
		return
	}

	// get the values from the message
	federationContextId := msg["federation_context_id"].(string)
	appPkgId := msg["app_pkg_id"].(string)
	vimId := msg["vim_id"].(string)
	config := msg["config"].(string)

	// get the federation context
	federation, err := f.services.FederationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		return
	}

	// check if the artefact exists in the given federation context
	artefact, err := f.services.ArtefactService.GetArtefactByAppPkgId(federationContextId, appPkgId)
	if err != nil {
		log.Printf("Error getting artefact: %v", err)
		return
	}

	// get the partner zone from the vimId
	var zoneId string
	for _, zone := range federation.PartnerOP.OfferedAvailabilityZones {
		if zone.VimId == vimId {
			zoneId = zone.ZoneId
			break
		}
	}

	if zoneId == "" {
		log.Printf("No zone found for vimId: %s", vimId)
		return
	}

	// make request to send to partner
	var request dto.InstantiateApplicationRequest
	request.TransactionId = uuid.New().String()
	request.AppId = artefact.Id
	request.AppProviderId = artefact.AppProviderId
	request.AppVersion = artefact.VersionInfo
	request.ZoneInfo = zoneId
	request.Config = config
	request.AppInstCallbackLink = "callback.link"

	// send request to partner
	appInstanceId, err := f.sendAppInstanceRequestToPartner(&federation, &request)
	if err != nil {
		log.Printf("Error sending app instance request to partner: %v", err)
		return
	}

	// make AppInstance object
	appInstance := models.AppInstance{
		Id:                  appInstanceId,
		FederationContextId: federationContextId,
		ArtefactId:          artefact.Id,
		Name:                "federated-instance",
		Description:         "not-local",
	}

	// save appinstance to database
	err = f.services.AppInstanceService.RegisterAppInstance(appInstance)
	if err != nil {
		log.Printf("Error saving app instance locally: %v", err)
		return
	}

	log.Printf("Successfully created app instance %s with partner operator and saved locally", appInstanceId)

	// send response to kafka
	_, err = f.services.KafkaClientService.Produce("responses", map[string]string{
		"msg_id":          msg["msg_id"].(string),
		"app_instance_id": appInstance.Id,
		"status":          "201",
	})
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}

}

func (f *FederationAppiNewCallback) sendAppInstanceRequestToPartner(federation *models.Federation, request *dto.InstantiateApplicationRequest) (string, error) {
	// use the stored access token from the federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// construct the partner's app instance endpoint URL
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/application/lcm", federation.FederationEndpoint, federation.PartnerOP.FederationContextId)

	// marshal the request
	payload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("partner returned error status %d", resp.StatusCode)
	}

	// parse response to get app instance ID
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	appInstanceId, ok := response["appInstanceId"].(string)
	if !ok {
		return "", fmt.Errorf("appInstanceId not found in partner response")
	}

	log.Printf("Partner response: app instance created with ID %s", appInstanceId)
	return appInstanceId, nil
}
