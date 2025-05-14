package ewbi

/*
 *
 * This file contains the implementation of the Federation Management Controller over the E/WBI
 * Exposes FederationManagement functionalities to other federators.
 *
 */

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type FederationManagementController struct {
	federationService *services.FederationService
}

// NewFederationManagementController creates a new instance of the FederationManagementController
func NewFederationManagementController(federationService *services.FederationService) *FederationManagementController {
	return &FederationManagementController{
		federationService: federationService,
	}
}

// @Summary Accept Federation Relationship
// @Description Establishes a new federation relationship with another federator with the provided data
// @Tags EWBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationRequestData body models.FederationRequestData true "Federation Request Data"
// @Success 200 {object} models.FederationResponseData
// @Failure 400 {object} models.ProblemDetails
// @Failure 500 {object} models.ProblemDetails
// @Router /ewbi/partner [post]
func (fmc *FederationManagementController) CreateFederationController(c *gin.Context) {
	log.Print("CreateFederationController - Establishing new federation relationship")

	// Check if the request data is valid
	var federationRequestData models.FederationRequestData
	if err := c.ShouldBindJSON(&federationRequestData); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body, missing parameters or wrong data type")
		return
	}

	log.Print("CreateFederationController - Creating federation object")
	federationResponseData := models.FederationResponseData{
		FederationContextId:          uuid.New().String(),
		PlatformCaps:                 &[]string{"MEC"},
		PartnerOPCountryCode:         "443",
		EdgeDiscoveryServiceEndPoint: &models.ServiceEndpoint{Fqdn: "edge-discovery-service.com", Port: 443},
		LcmServiceEndPoint:           &models.ServiceEndpoint{Fqdn: "lcm-service.com", Port: 443},
		FederationExpiryDate:         time.Now().AddDate(1, 0, 0),
		FederationRenewalDate:        time.Now().AddDate(0, 6, 0),
	}

	healthInfo := models.FederationHealthInfo{
		FederationStatus:    &models.State{AlarmState: "CLEAR"},
		FederationStartTime: time.Now(),
		NumOfAcceptedZones:  "0",
		NumOfActiveAlarms:   "0",
		NumOfApplications:   "0",
	}

	federation := models.Federation{
		PartnerOP:     federationResponseData,
		OriginOP:      federationRequestData,
		HealthInfo:    healthInfo,
		IsEstablished: true,
		IsOriginOP:    false,
	}

	log.Print("CreateFederationController - Federation object created, storing in database")
	// Store the federation object in the database
	federation, err := fmc.federationService.CreateFederation(federation)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error creating the Federation object")
		return
	}

	log.Print("CreateFederationController - Federation relationship established successfully")
	c.JSON(http.StatusOK, federation.PartnerOP)
}

// @Summary Remove Federation Relationship
// @Description Removes a federation relationship by its federationContextId
// @Tags EWBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} map[string]string "status: Federation removed successfully"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/partner [delete]
func (fmc *FederationManagementController) RemoveFederationController(c *gin.Context) {
	log.Print("RemoveFederationRelationshipController - Removing federation relationship")

	// Delete the federation object from the database
	federationContextId := c.Param("federationContextId")
	err := fmc.federationService.DeleteFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error removing the Federation object from database")
		return
	}

	log.Print("RemoveFederationRelationshipController - Federation removed successfully")
	c.JSON(http.StatusOK, gin.H{"status": "Federation removed successfully"})
}

// @Summary Get Federation Meta Info
// @Description Retrieves metadata information about a federation based on the federationContextId
// @Tags EWBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} models.FederationMetaInfo
// @Failure 400 {object} models.ProblemDetails "Invalid Federation Context ID"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/partner [get]
func (fmc *FederationManagementController) GetFederationMetaInfoController(c *gin.Context) {
	log.Print("GetFederationMetaInfoController - Retrieving federation meta information")

	// Get federation details from the database
	federationContextId := c.Param("federationContextId")
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving federation information")
		return
	}

	response := models.FederationMetaInfo{
		EdgeDiscoveryServiceEndPoint: federation.PartnerOP.EdgeDiscoveryServiceEndPoint,
		OfferedAvailabilityZones:     federation.PartnerOP.OfferedAvailabilityZones,
		AllowedMobileNetworkIds:      federation.OriginOP.OrigOPMobileNetworkCodes,
		AllowedFixedNetworkIds:       federation.OriginOP.OrigOPFixedNetworkCodes,
		LcmServiceEndPoint:           federation.PartnerOP.LcmServiceEndPoint,
		PlatformCaps:                 federation.PartnerOP.PlatformCaps,
	}

	log.Print("GetFederationMetaInfoController - Federation meta information retrieved successfully")
	c.JSON(http.StatusOK, response)
}

// @Summary Update Federation Details
// @Description Updates a federation object with the given federationContextId and patch parameters
// @Tags EWBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param patchParams body models.FederationPatchParams true "Patch Parameters"
// @Success 200 {object} map[string]string "status: Federation updated successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request or federation not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /ewbi/{federationContextId}/partner [patch]
func (fmc *FederationManagementController) UpdateFederationController(c *gin.Context) {
	log.Print("UpdateFederationController - Updating federation, checking patch parameters from body")
	var patchParams models.FederationPatchParams
	if err := c.ShouldBindJSON(&patchParams); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body, missing parameters or wrong data type")
		return
	}

	federationContextId := c.Param("federationContextId")

	log.Print("UpdateFederationController - Checking validity of patch parameters")
	if patchParams.ObjectType != "MOBILE_NETWORK_CODES" && patchParams.ObjectType != "FIXED_NETWORK_CODES" {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid ObjectType, must be either 'MOBILE_NETWORK_CODES' or 'FIXED_NETWORK_CODES'")
		return
	}

	if patchParams.OperationType != "ADD_CODES" && patchParams.OperationType != "REMOVE_CODES" && patchParams.OperationType != "UPDATE_CODES" {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid OperationType, must be either 'ADD_CODES', 'REMOVE_CODES' or 'UPDATE_CODES'")
		return
	}

	log.Print("UpdateFederationController - Updating federation object")
	err := fmc.federationService.PatchFederation(federationContextId, patchParams)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error updating the Federation object: "+err.Error())
		return
	}

	log.Print("UpdateFederationController - Federation updated successfully")
	c.JSON(http.StatusOK, gin.H{"status": "Federation updated successfully"})
}

// @Summary Get Federation Context Identifier
// @Description Retrieves the federationContextId using the accessToken
// @Tags EWBI - FederationManagement
// @Accept json
// @Produce json
// @Param accessToken query string true "Access Token"
// @Success 200 {object} map[string]string "federationContextId: id"
// @Failure 400 {object} models.ProblemDetails "No Federation found with the given accessToken"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/fed-context-id [get]
func (fmc *FederationManagementController) GetFederationContextIdentifierController(c *gin.Context) {
	log.Print("GetFederationContextIdentifier - Retrieving federationContextId using the accessToken")

	// Check if the a federation exists with the given accessToken
	accessToken := c.Query("accessToken")
	if !fmc.federationService.ExistsFederationWithAccessToken(accessToken) {
		utils.HandleProblem(c, http.StatusBadRequest, "No Federation found with the given accessToken")
		return
	}

	// Retrieve the federation through the accessToken
	federation, err := fmc.federationService.GetFederationByAccessToken(accessToken)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving federation information: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"federationContextId": federation.PartnerOP.FederationContextId})
}

// @Summary Get Local Federation Health Info
// @Description Checks the health status of the federation
// @Tags EWBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} models.FederationHealthInfo
// @Failure 400 {object} models.ProblemDetails "Invalid federationContextId"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/health [get]
func (fmc *FederationManagementController) GetFederationHealthController(c *gin.Context) {
	log.Print("HealthCheckController - Checking health status of the federator")

	log.Print("HealthCheckController - Checking if federation exists")
	// Check if the federation exists
	federationContextId := c.Param("federationContextId")
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		utils.HandleProblem(c, http.StatusBadRequest, "No Federation found with the given federationContextId")
		return
	}

	log.Print("HealthCheckController - Federation exists, retrieving health information")
	// Retrieve the federation information
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving federation information")
		return
	}

	log.Print("HealthCheckController - Returning health information")
	c.JSON(http.StatusOK, federation.HealthInfo)
}

// @Summary Renew Federation Relationship
// @Description Renews the federation relationship with the given federationContextId
// @Tags EWBI - FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} map[string]string "federationContextId: id, federationRenewalDate: time, federationExpiryDate: time"
// @Failure 400 {object} models.ProblemDetails "Invalid federationContextId"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/renew [post]
func (fmc *FederationManagementController) RenewFederationController(c *gin.Context) {
	log.Print("RenewFederationController - Renewing federation relationship")

	federationContextId := c.Param("federationContextId")
	log.Print("RenewFederationController - Checking if federation exists")
	// Check if the federation exists
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		utils.HandleProblem(c, http.StatusBadRequest, "No Federation found with the given federationContextId")
		return
	}

	log.Print("RenewFederationController - Getting federation object")
	// Get the federation object
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving the Federation object")
		return
	}

	log.Print("RenewFederationController - Renewing federation relationship")
	// Renew the federation relationship
	federation, err = fmc.federationService.RenewFederation(federation)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error renewing the Federation object")
		return
	}

	response := dto.FederationRenewalResponseData{
		FederationContextId:   federation.PartnerOP.FederationContextId,
		FederationRenewalDate: federation.PartnerOP.FederationRenewalDate,
		FederationExpiryDate:  federation.PartnerOP.FederationExpiryDate,
	}

	log.Print("RenewFederationController - Federation relationship renewed successfully")
	c.JSON(http.StatusOK, response)
}
