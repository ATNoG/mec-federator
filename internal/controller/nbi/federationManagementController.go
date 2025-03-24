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

// @Summary Initiate Federation Relationship
// @Description Initiates the federation establishment procedure with another federator
// @Tags NBI - FederationManagement
// @Accept  json
// @Produce  json
// @Param   federationRequest  body  models.FederationInitiateRequest  true  "Federation and Auth Endpoints"
// @Success 200 {object} models.FederationResponseData
// @Failure 400 {object} models.ProblemDetails "Invalid request body or missing fields"
// @Failure 500 {object} models.ProblemDetails "Internal error during federation process"
// @Router /nbi/federation/v1/partner [post]
func (fmc *FederationManagementController) InitiateFederationController(c *gin.Context) {
	log.Print("InitiateFederationController - Initiating federation")

	// get federation and auth endpoints from request body
	var requestBody models.FederationInitiateRequest
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
	err := httpClient.HttpRequestWithAuthAndUnmarshal(c, http.MethodPost, "/federation/v1/partner", federationRequestData, &federationResponseData)
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

// @Summary Remove Federation Relationship
// @Description Removes a previously established federation relationship by its federationContextId
// @Tags NBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} map[string]string "status: Federation removed successfully"
// @Failure 400 {object} models.ProblemDetails "federationContextId must be provided"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /nbi/federation/v1/{federationContextId}/partner [delete]
func (fmc *FederationManagementController) RemoveFederationController(c *gin.Context) {
	log.Print("RemoveFederationController - Removing federation")

	// get federation context id from request
	federationContextId := c.Param("federationContextId")

	if federationContextId == "" {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "federationContextId must be provided"
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	log.Print("RemoveFederationController - Getting HttpClient for Federation")

	// get http client for federation
	httpClient, err := fmc.federationHttpClientManager.Get(federationContextId)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "No http client found with the given federationContextId"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("RemoveFederationController - Sending Federation Removal Request to Partner")

	// send federation removal request to other federator
	resp, err := httpClient.HttpRequestWithAuth(
		c,
		http.MethodDelete,
		"/federation/v1/"+federationContextId+"/partner",
		nil)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = err.Error()
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	if resp.StatusCode != http.StatusOK {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Error removing federation from partner"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	// remove federation from database
	fmc.federationService.DeleteFederation(federationContextId)

	log.Print("RemoveFederationController - Federation removed successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Federation removed successfully"})
}

// @Summary Get Federation Meta Info from partner
// @Description Retrieves federation information from the partner federator
// @Tags NBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} models.FederationHealthInfo
// @Failure 400 {object} models.ProblemDetails "federationContextId must be provided"
// @Failure 404 {object} models.ProblemDetails "Federation not found"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /nbi/federation/v1/{federationContextId}/partner [get]
func (fmc *FederationManagementController) GetFederationMetaInfoController(c *gin.Context) {
	log.Print("GetFederationMetaInfoController - Getting federation meta info from partner")

	// get federation context id from request
	federationContextId := c.Param("federationContextId")
	if federationContextId == "" {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "federationContextId must be provided"
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	log.Print("GetFederationMetaInfoController - Checking if federation exists")

	// check if federation exists
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		problemDetails := utils.NewProblemDetails(http.StatusNotFound)
		problemDetails.Detail = "Federation not found"
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}

	log.Print("GetFederationMetaInfoController - Getting HttpClient for Federation")

	// get http client for federation
	httpClient, err := fmc.federationHttpClientManager.Get(federationContextId)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "No http client found with the given federationContextId"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("GetFederationMetaInfoController - Sending meta info request to partner")

	// send meta info request to other federator
	var metaInfo models.FederationMetaInfo
	err = httpClient.HttpRequestWithAuthAndUnmarshal(
		c,
		http.MethodGet,
		"/federation/v1/"+federationContextId+"/partner",
		nil,
		&metaInfo)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = err.Error()
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("GetFederationMetaInfoController - Federation meta info retrieved successfully")
}

// @Summary Get Federation Health Info
// @Description Retrieves the health information of the federation partner
// @Tags NBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} models.FederationHealthInfo
// @Failure 400 {object} models.ProblemDetails "federationContextId must be provided"
// @Failure 404 {object} models.ProblemDetails "Federation not found"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /nbi/federation/v1/{federationContextId}/health [get]
func (fmc *FederationManagementController) GetFederationHealthController(c *gin.Context) {
	log.Print("GetFederationHealthController - Getting federation health info")

	// get federation context id from request
	federationContextId := c.Param("federationContextId")

	if federationContextId == "" {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "federationContextId must be provided"
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	log.Print("GetFederationHealthController - Checking if federation exists")

	// check if federation exists
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		problemDetails := utils.NewProblemDetails(http.StatusNotFound)
		problemDetails.Detail = "Federation not found"
		c.JSON(http.StatusNotFound, problemDetails)
		return
	}

	log.Print("GetFederationHealthController - Getting HttpClient for Federation")

	// get http client for federation
	httpClient, err := fmc.federationHttpClientManager.Get(federationContextId)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "No http client found with the given federationContextId"
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("GetFederationHealthController - Sending health check request to partner")

	// send health check request to other federator
	var healthInfo models.FederationHealthInfo
	err = httpClient.HttpRequestWithAuthAndUnmarshal(
		c,
		http.MethodGet,
		"/federation/v1/"+federationContextId+"/health",
		nil,
		&healthInfo)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = err.Error()
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("GetFederationHealthController - Federation health info retrieved successfully")

	c.JSON(http.StatusOK, healthInfo)
}
