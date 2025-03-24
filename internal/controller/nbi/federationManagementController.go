package nbi

/*
 *
 * This file contains the implementation of the Federation Management Controller over the NBI.
 * Exposes Federator functionalities to the orchestrators it manages.
 *
 */

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type FederationManagementController struct {
	federationService           *services.FederationService
	federationHttpClientManager *services.FederationHttpClientManager
}

// NewFederationManagementController creates a new instance of the FederationManagementController
func NewFederationManagementController(fs *services.FederationService, hcm *services.FederationHttpClientManager) *FederationManagementController {
	return &FederationManagementController{
		federationHttpClientManager: hcm,
		federationService:           fs,
	}
}

func (fmc *FederationManagementController) InitiateFederationController(c *gin.Context) {
	log.Print("InitiateFederationController - Initiating federation")

	// get federation and auth endpoints from request body
	var requestBody struct {
		FederationEndpoint string `json:"federationEndpoint" binding:"required"`
		AuthEndpoint       string `json:"authEndpoint" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "Invalid request body: " + err.Error()
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	federationEndpoint := requestBody.FederationEndpoint
	authEndpoint := requestBody.AuthEndpoint

	if federationEndpoint == "" || authEndpoint == "" {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "federationEndpoint and authEndpoint must be provided"
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	log.Print("InitiateFederationController - Creating Federation Request Data")

	// create federation request data
	federationRequestData := models.FederationRequestData{
		OrigOPCountryCode:       "351",
		OrigOPFixedNetworkCodes: &[]string{"351", "352"},
		InitialDate:             time.Now(),
		PartnerStatusLink:       "https://status.link",
	}

	log.Print("InitiateFederationController - Creating HTTP Client for Federation")

	// make http client for federation
	httpClientConfig := services.HttpClientConfig{
		BaseUrl:       federationEndpoint,
		TokenEndpoint: authEndpoint,
		ClientId:      config.AppConfig.OAuth2ClientId,
		ClientSecret:  config.AppConfig.OAuth2ClientSecret,
	}

	httpClient := services.NewHttpClient(httpClientConfig)

	log.Print("InitiateFederationController - Sending Federation Request to Partner")

	// send federation request to other federator and unmarshal response
	var federationResponseData models.FederationResponseData
	err := httpClient.PostJSONWithAuthAndUnmarshal(c, federationEndpoint+"/federation/v1/partner", federationRequestData, &federationResponseData)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = err.Error()
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("InitiateFederationController - Creating Federation Object")

	// create federation object
	federation := models.Federation{
		PartnerOP:          federationResponseData,
		OriginOP:           federationRequestData,
		IsEstablished:      true,
		FederationEndpoint: federationEndpoint,
		AuthEndpoint:       authEndpoint,
		IsOriginOP:         true,
	}

	log.Print("InitiateFederationController - Saving Federation to Database")

	// save federation to database
	fmc.federationService.CreateFederation(federation)

	log.Print("InitiateFederationController - Registering HttpClient for Federation")

	// register httpClient for federation
	fmc.federationHttpClientManager.Register(federationResponseData.FederationContextId, httpClient)

	log.Print("InitiateFederationController - Federation initiated successfully")

	c.JSON(http.StatusOK, federationResponseData)
}
