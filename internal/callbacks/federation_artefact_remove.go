package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/services"
)

type FederationArtefactRemoveCallback struct {
	authService        *services.AuthService
	httpClientService  *services.HttpClientService
	kafkaClientService *services.KafkaClientService
	federationService  *services.FederationService
	artefactService    *services.ArtefactService
}

func NewFederationArtefactRemoveCallback(authService *services.AuthService, httpClientService *services.HttpClientService, kafkaClientService *services.KafkaClientService, federationService *services.FederationService, artefactService *services.ArtefactService) *FederationArtefactRemoveCallback {
	return &FederationArtefactRemoveCallback{authService: authService, httpClientService: httpClientService, kafkaClientService: kafkaClientService, federationService: federationService, artefactService: artefactService}
}

// receives info about an artefact to remove from a certain federation
func (f *FederationArtefactRemoveCallback) HandleMessage(message *sarama.ConsumerMessage) {
	// unmarshal the message
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling response message: %v", err)
		return
	}

	appPkgId := msg["app_pkg_id"].(string)
	federationContextId := msg["federation_context_id"].(string)

	// get the federation context
	federation, err := f.federationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("Error getting federation: %v", err)
		return
	}

	// check if an artefact with this appPkgId exists in the federation
	artefact, err := f.artefactService.GetArtefactByAppPkgId(federationContextId, appPkgId)
	if err != nil {
		log.Printf("No artefact found with appPkgId %s in federation %s: %v", appPkgId, federationContextId, err)
		return
	}

	// remove the artefact from the partner operator
	err = f.removeArtefactFromPartner(&federation, artefact.Id)
	if err != nil {
		log.Printf("Error removing artefact from partner: %v", err)
		return
	}

	// remove the artefact locally
	err = f.artefactService.RemoveArtefact(federationContextId, artefact.Id)
	if err != nil {
		log.Printf("Error removing artefact locally: %v", err)
		return
	}

	log.Printf("Successfully removed artefact %s from partner operator and locally", appPkgId)
}

func (f *FederationArtefactRemoveCallback) removeArtefactFromPartner(federation *models.Federation, artefactId string) error {
	// use the stored access token from the federation
	accessToken := federation.OriginOP.AccessToken.AccessToken

	// construct the partner's artefact endpoint URL for deletion
	partnerEndpoint := fmt.Sprintf("%s/federation/v1/ewbi/%s/artefact/%s", federation.FederationEndpoint, federation.PartnerOP.FederationContextId, artefactId)

	// create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	auth := services.NewBearerTokenAuth(accessToken)

	// make the HTTP DELETE request
	resp, err := f.httpClientService.DoRequest(
		ctx,
		http.MethodDelete,
		partnerEndpoint,
		nil,
		headers,
		auth)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("partner returned error status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Partner response: %s", string(respBody))
	return nil
}
