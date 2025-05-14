package nbi

/*
 *
 * This file contains the implementation of the Federation Management Controller over the NBI.
 * Exposes Federator functionalities to the orchestrators it manages.
 *
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type FederationManagementController struct {
	federationService *services.FederationService
	httpClientService *services.HttpClientService
}

// NewFederationManagementController creates a new instance of the FederationManagementController
func NewFederationManagementController(fs *services.FederationService, hcp *services.HttpClientService) *FederationManagementController {
	return &FederationManagementController{
		federationService: fs,
		httpClientService: hcp,
	}
}

// @Summary Initiate Federation Relationship with a partner OP
// @Description Initiates the federation establishment procedure with another federator
// @Tags NBI - FederationManagement
// @Accept  json
// @Produce  json
// @Param   federationRequest  body  models.FederationInitiateRequestData  true  "Federation and Auth Endpoints"
// @Success 200 {object} models.FederationResponseData
// @Failure 400 {object} models.ProblemDetails "Invalid request body or missing fields"
// @Failure 500 {object} models.ProblemDetails "Internal error during federation process"
// @Router /nbi/partner [post]
func (fmc *FederationManagementController) InitiateFederationController(c *gin.Context) {
	log.Print("InitiateFederationController - Initiating Federation establishment procedure")

	// Bind request body to struct
	var requestBody dto.FederationInitiateRequestData
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Check if required fields are present
	if requestBody.FederationEndpoint == "" {
		utils.HandleProblem(c, http.StatusBadRequest, "Missing FederationEndpoint in request body")
		return
	}
	if requestBody.AuthEndpoint == "" {
		utils.HandleProblem(c, http.StatusBadRequest, "Missing AuthEndpoint in request body")
		return
	}

	log.Print("InitiateFederationController - Fetching Access Token from Auth Endpoint")
	accessTokenUrl := requestBody.AuthEndpoint
	headers := map[string]string{"Content-Type": "application/json"}
	accessTokenRequestData := dto.AccessTokenRequestData{
		ClientId:     config.AppConfig.OAuth2ClientId,
		ClientSecret: config.AppConfig.OAuth2ClientSecret,
	}
	var accessToken models.AccessToken

	// Marshal request data to json
	jsonPayload, err := json.Marshal(accessTokenRequestData)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Send request to auth endpoint
	resp, err := fmc.httpClientService.DoRequest(
		c,
		http.MethodPost,
		accessTokenUrl,
		bytes.NewBuffer(jsonPayload),
		headers,
		nil)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	// Check if response is OK
	if resp.StatusCode != http.StatusOK {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error fetching access token from auth endpoint")
		return
	}

	// Decode response body to struct
	err = json.NewDecoder(resp.Body).Decode(&accessToken)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Print("InitiateFederationController - Sending Federation Request to Partner")
	createFederationUrl := fmt.Sprintf("%s%s", requestBody.FederationEndpoint, "/federation/v1/ewbi/partner")
	authStrat := services.NewBearerTokenAuth(accessToken.AccessToken)
	federationRequestData := models.FederationRequestData{
		OrigOPCountryCode:       "351",
		OrigOPFixedNetworkCodes: &[]string{"351", "352"},
		InitialDate:             time.Now(),
		PartnerStatusLink:       "https://status.link",
		AccessToken:             accessToken,
	}
	var federationResponseData models.FederationResponseData

	// Marshal request data to json
	jsonPayload, err = json.Marshal(federationRequestData)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	// Send request to partner federator
	resp, err = fmc.httpClientService.DoRequest(
		c,
		http.MethodPost,
		createFederationUrl,
		bytes.NewBuffer(jsonPayload),
		headers,
		authStrat)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError, "Error sending federation request to partner")
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}
	defer resp.Body.Close()

	// Check if response is OK
	if resp.StatusCode != http.StatusOK {
		utils.HandleProblem(c, http.StatusInternalServerError, "Partner federator returned an error")
		return
	}

	// Decode response body to struct
	err = json.NewDecoder(resp.Body).Decode(&federationResponseData)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Print("InitiateFederationController - Creating Federation Object")
	federation := models.Federation{
		PartnerOP:          federationResponseData,
		OriginOP:           federationRequestData,
		IsEstablished:      true,
		FederationEndpoint: requestBody.FederationEndpoint,
		AuthEndpoint:       requestBody.AuthEndpoint,
		IsOriginOP:         true,
	}

	log.Print("InitiateFederationController - Saving Federation to Database")
	// Save federation to database
	_, err = fmc.federationService.CreateFederation(federation)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	log.Print("InitiateFederationController - Federation initiated successfully")
	c.JSON(http.StatusOK, federationResponseData)
}

// @Summary Remove Federation Relationship with a partner OP
// @Description Removes a previously established federation relationship by its federationContextId
// @Tags NBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} map[string]string "status: Federation removed successfully"
// @Failure 400 {object} models.ProblemDetails "federationContextId must be provided"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /nbi/{federationContextId}/partner [delete]
func (fmc *FederationManagementController) RemoveFederationController(c *gin.Context) {
	log.Print("RemoveFederationController - Removing federation")

	federationContextId := c.Param("federationContextId")
	if federationContextId == "" {
		utils.HandleProblem(c, http.StatusBadRequest, "federationContextId must be provided")
		return
	}

	log.Print("RemoveFederationController - Checking if federation exists")
	// Check if federation exists
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		utils.HandleProblem(c, http.StatusNotFound, "Federation not found")
		return
	}

	log.Print("RemoveFederationController - Getting Federation Url from Database")
	// Get federation details from database
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting federation details from database:"+err.Error())
		return
	}

	log.Print("RemoveFederationController - Sending Federation Removal Request to Partner")
	deleteFederationUrl := fmt.Sprintf("%s%s%s%s", federation.FederationEndpoint, "/federation/v1/ewbi/", federationContextId, "/partner")
	authStrat := services.NewBearerTokenAuth(federation.OriginOP.AccessToken.AccessToken)

	// Send request to partner federator
	resp, err := fmc.httpClientService.DoRequest(
		c,
		http.MethodDelete,
		deleteFederationUrl,
		nil,
		nil,
		authStrat)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error sending request to partner: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Check if response is OK
	if resp.StatusCode != http.StatusOK {
		utils.HandleProblem(c, http.StatusInternalServerError, "Partner federator returned an error")
		return
	}

	log.Print("RemoveFederationController - Removing Federation from Database")
	// Remove federation from database
	err = fmc.federationService.DeleteFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error deleting federation from database: "+err.Error())
		return
	}

	log.Print("RemoveFederationController - Federation removed successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Federation removed successfully"})
}

// @Summary Get Federation Meta Info from partner OP
// @Description Retrieves federation information from the partner federator
// @Tags NBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} models.FederationHealthInfo
// @Failure 400 {object} models.ProblemDetails "federationContextId must be provided"
// @Failure 404 {object} models.ProblemDetails "Federation not found"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /nbi/{federationContextId}/partner [get]
func (fmc *FederationManagementController) GetFederationMetaInfoController(c *gin.Context) {
	log.Print("GetFederationMetaInfoController - Getting federation meta info from partner")

	federationContextId := c.Param("federationContextId")
	if federationContextId == "" {
		utils.HandleProblem(c, http.StatusBadRequest, "federationContextId must be provided")
		return
	}

	log.Print("GetFederationMetaInfoController - Checking if federation exists")
	// Check if federation exists
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		utils.HandleProblem(c, http.StatusNotFound, "Federation not found")
		return
	}

	log.Print("GetFederationMetaInfoController - Getting Federation Url from Database")
	// Get federation details from database
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting federation details from database:"+err.Error())
		return
	}

	log.Print("GetFederationMetaInfoController - Sending meta info request to partner")
	getMetaInfoUrl := fmt.Sprintf("%s%s%s%s", federation.FederationEndpoint, "/federation/v1/ewbi/", federationContextId, "/partner")
	authStrat := services.NewBearerTokenAuth(federation.OriginOP.AccessToken.AccessToken)
	var metaInfo models.FederationMetaInfo

	// Send meta info request to partner federator
	resp, err := fmc.httpClientService.DoRequest(
		c,
		http.MethodGet,
		getMetaInfoUrl,
		nil,
		nil,
		authStrat)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error sending request to partner: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Check if response is OK
	if resp.StatusCode != http.StatusOK {
		utils.HandleProblem(c, http.StatusInternalServerError, "Partner federator returned error "+fmt.Sprintf("%d", resp.StatusCode))
		return
	}

	// Decode response body to struct
	err = json.NewDecoder(resp.Body).Decode(&metaInfo)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error decoding response body: "+err.Error())
		return
	}

	log.Print("GetFederationMetaInfoController - Federation meta info retrieved successfully")
	c.JSON(http.StatusOK, metaInfo)
}

// @Summary Get Federation Health Info from partner OP
// @Description Retrieves the health information of the federation partner
// @Tags NBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} models.FederationHealthInfo
// @Failure 400 {object} models.ProblemDetails "federationContextId must be provided"
// @Failure 404 {object} models.ProblemDetails "Federation not found"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /nbi/{federationContextId}/health [get]
func (fmc *FederationManagementController) GetFederationHealthController(c *gin.Context) {
	log.Print("GetFederationHealthController - Getting federation health info")

	federationContextId := c.Param("federationContextId")
	if federationContextId == "" {
		utils.HandleProblem(c, http.StatusBadRequest, "federationContextId must be provided")
		return
	}

	// Check if federation exists
	log.Print("GetFederationHealthController - Checking if federation exists")
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		utils.HandleProblem(c, http.StatusNotFound, "Federation not found")
		return
	}

	// Get federation details from database
	log.Print("GetFederationHealthController - Getting Federation Url from Database")
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting federation details from database: "+err.Error())
		return
	}

	log.Print("GetFederationHealthController - Sending health check request to partner")
	getHealthUrl := fmt.Sprintf("%s%s%s%s", federation.FederationEndpoint, "/federation/v1/ewbi/", federationContextId, "/health")
	authStrat := services.NewBearerTokenAuth(federation.OriginOP.AccessToken.AccessToken)
	var healthInfo models.FederationHealthInfo

	// Send request to partner federator
	resp, err := fmc.httpClientService.DoRequest(
		c,
		http.MethodGet,
		getHealthUrl,
		nil,
		nil,
		authStrat)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error sending request to partner: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Check if response is OK
	if resp.StatusCode != http.StatusOK {
		utils.HandleProblem(c, http.StatusInternalServerError, "Partner federator returned an error")
		return
	}

	// Decode response body to struct
	err = json.NewDecoder(resp.Body).Decode(&healthInfo)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error decoding response body: "+err.Error())
		return
	}

	log.Print("GetFederationHealthController - Federation health info retrieved successfully")
	c.JSON(http.StatusOK, healthInfo)
}

func (fmc *FederationManagementController) FetchFedContextIdFromPartnerController(c *gin.Context) {

}

// @Summary Request Federation Renewal with partner OP
// @Description Requests the renewal of the federation with the partner federator
// @Tags NBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} models.FederationRenewalResponseData
// @Failure 400 {object} models.ProblemDetails "federationContextId must be provided"
// @Failure 404 {object} models.ProblemDetails "Federation not found"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /nbi/{federationContextId}/renew [post]
func (fmc *FederationManagementController) RequestFederationRenewalController(c *gin.Context) {
	log.Print("RequestFederationRenewalController - Requesting federation renewal")

	federationContextId := c.Param("federationContextId")
	if federationContextId == "" {
		utils.HandleProblem(c, http.StatusBadRequest, "federationContextId must be provided")
		return
	}

	// Check if federation exists
	log.Print("RequestFederationRenewalController - Checking if federation exists")
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		utils.HandleProblem(c, http.StatusNotFound, "Federation not found")
		return
	}

	// Get federation details from database
	log.Print("RequestFederationRenewalController - Getting Federation Url from Database")
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting federation details from database: "+err.Error())
		return
	}

	log.Print("RequestFederationRenewalController - Sending renewal request to partner")
	renewFederationUrl := fmt.Sprintf("%s%s%s%s", federation.FederationEndpoint, "/federation/v1/ewbi/", federationContextId, "/renew")
	authStrat := services.NewBearerTokenAuth(federation.OriginOP.AccessToken.AccessToken)
	var renewalResponse dto.FederationRenewalResponseData

	// Send request to partner federator
	resp, err := fmc.httpClientService.DoRequest(
		c,
		http.MethodPost,
		renewFederationUrl,
		nil,
		nil,
		authStrat)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error sending request to partner: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Check if response is OK
	if resp.StatusCode != http.StatusOK {
		utils.HandleProblem(c, http.StatusInternalServerError, "Partner federator returned an error")
		return
	}

	// Decode response body to struct
	err = json.NewDecoder(resp.Body).Decode(&renewalResponse)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error deconding responde body: "+err.Error())
		return
	}

	log.Print("RequestFederationRenewalController - Federation renewal successfully by partner")

	federation.PartnerOP.FederationRenewalDate = renewalResponse.FederationRenewalDate
	federation.PartnerOP.FederationExpiryDate = renewalResponse.FederationExpiryDate

	log.Print("RequestFederationRenewalController - Updating Federation in Database")
	// Update federation in database
	err = fmc.federationService.UpdateFederation(federation)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Could not update federation in database: "+err.Error())
		return
	}

	log.Print("RequestFederationRenewalController - Federation renewed successfully")
	c.JSON(http.StatusOK, renewalResponse)
}
