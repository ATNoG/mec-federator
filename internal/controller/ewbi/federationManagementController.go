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
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type FederationManagementController struct {
	federationService  *services.FederationService
	zoneService        *services.ZoneService
	artefactService    *services.ArtefactService
	appInstanceService *services.AppInstanceService
}

// NewFederationManagementController creates a new instance of the FederationManagementController
func NewFederationManagementController(federationService *services.FederationService, zoneService *services.ZoneService, artefactService *services.ArtefactService, appInstanceService *services.AppInstanceService) *FederationManagementController {
	return &FederationManagementController{
		federationService:  federationService,
		zoneService:        zoneService,
		artefactService:    artefactService,
		appInstanceService: appInstanceService,
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
	log.Print("CreateFederationController - Starting federation creation process")
	
	// Check if the request data is valid
	var federationRequestData models.FederationRequestData
	if err := c.ShouldBindJSON(&federationRequestData); err != nil {
		log.Printf("CreateFederationController - Invalid request body: %v", err)
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body, missing parameters or wrong data type")
		return
	}
	
	log.Printf("CreateFederationController - Received federation request from partner: %s", federationRequestData.OrigOPFederationId)

	// Create the federation response data
	log.Print("CreateFederationController - Creating federation response data")
	federationResponseData := models.FederationResponseData{
		PartnerOPFederationId:        config.AppConfig.OperatorId,
		FederationContextId:          uuid.New().String(),
		PlatformCaps:                 &[]string{"MEC"},
		PartnerOPCountryCode:         "443",
		EdgeDiscoveryServiceEndPoint: &models.ServiceEndpoint{Fqdn: "edge-discovery-service.com", Port: 443},
		LcmServiceEndPoint:           &models.ServiceEndpoint{Fqdn: "lcm-service.com", Port: 443},
		FederationExpiryDate:         time.Now().AddDate(1, 0, 0),
		FederationRenewalDate:        time.Now().AddDate(0, 6, 0),
	}

	// Get the local available zones
	log.Print("CreateFederationController - Retrieving local zones")
	localZones, err := fmc.zoneService.GetLocalZones()
	if err != nil {
		log.Printf("CreateFederationController - Error retrieving local zones: %v", err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving local zones")
		return
	}
	
	log.Printf("CreateFederationController - Retrieved %d local zones", len(localZones))

	// Add zones to the federation response data
	federationResponseData.OfferedAvailabilityZones = localZones

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

	// Store the federation object in the database
	log.Printf("CreateFederationController - Storing federation in database with contextId: %s", federation.PartnerOP.FederationContextId)
	federation, err = fmc.federationService.CreateFederation(federation)
	if err != nil {
		log.Printf("CreateFederationController - Error creating federation in database: %v", err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error creating the Federation object")
		return
	}
	
	log.Printf("CreateFederationController - Federation created successfully with contextId: %s", federation.PartnerOP.FederationContextId)

	log.Printf("CreateFederationController - Returning federation response for contextId: %s", federation.PartnerOP.FederationContextId)
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
	// Get the federation context id
	federationContextId := c.Param("federationContextId")
	log.Printf("RemoveFederationController - Starting federation removal for contextId: %s", federationContextId)

	// Check if there are any app instances running in the federation
	log.Printf("RemoveFederationController - Checking for running app instances in federation: %s", federationContextId)
	appInstances, err := fmc.appInstanceService.GetAppInstancesByFederationContextId(federationContextId)
	if err != nil {
		log.Printf("RemoveFederationController - Error retrieving app instances: %v", err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving app instances")
		return
	}

	if len(appInstances) > 0 {
		log.Printf("RemoveFederationController - Cannot remove federation %s: %d app instances still running", federationContextId, len(appInstances))
		utils.HandleProblem(c, http.StatusBadRequest, "There are still app instances running in the federation")
		return
	}
	
	log.Printf("RemoveFederationController - No app instances found in federation: %s", federationContextId)

	// Check if there are any artefacts stored in the federation
	log.Printf("RemoveFederationController - Checking for artefacts in federation: %s", federationContextId)
	artefacts, err := fmc.artefactService.GetArtefactsByFederationContextId(federationContextId)
	if err != nil {
		log.Printf("RemoveFederationController - Error retrieving artefacts: %v", err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving artefacts")
		return
	}

	if len(artefacts) > 0 {
		log.Printf("RemoveFederationController - Cannot remove federation %s: %d artefacts still stored", federationContextId, len(artefacts))
		utils.HandleProblem(c, http.StatusBadRequest, "There are still artefacts stored in the federation")
		return
	}
	
	log.Printf("RemoveFederationController - No artefacts found in federation: %s", federationContextId)

	// Delete the federation object from the database
	log.Printf("RemoveFederationController - Deleting federation from database: %s", federationContextId)
	err = fmc.federationService.DeleteFederation(federationContextId)
	if err != nil {
		log.Printf("RemoveFederationController - Error deleting federation from database: %v", err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error removing the Federation object from database")
		return
	}
	
	log.Printf("RemoveFederationController - Federation %s removed successfully", federationContextId)

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
	// Get federation details from the database
	federationContextId := c.Param("federationContextId")
	log.Printf("GetFederationMetaInfoController - Retrieving meta info for federation: %s", federationContextId)
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
	federationContextId := c.Param("federationContextId")
	log.Printf("UpdateFederationController - Starting federation update for contextId: %s", federationContextId)
	
	var patchParams models.FederationPatchParams
	if err := c.ShouldBindJSON(&patchParams); err != nil {
		log.Printf("UpdateFederationController - Invalid request body for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body, missing parameters or wrong data type")
		return
	}
	
	log.Printf("UpdateFederationController - Patch parameters for federation %s: ObjectType=%s, OperationType=%s", federationContextId, patchParams.ObjectType, patchParams.OperationType)

	if patchParams.ObjectType != "MOBILE_NETWORK_CODES" && patchParams.ObjectType != "FIXED_NETWORK_CODES" {
		log.Printf("UpdateFederationController - Invalid ObjectType for federation %s: %s", federationContextId, patchParams.ObjectType)
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid ObjectType, must be either 'MOBILE_NETWORK_CODES' or 'FIXED_NETWORK_CODES'")
		return
	}

	if patchParams.OperationType != "ADD_CODES" && patchParams.OperationType != "REMOVE_CODES" && patchParams.OperationType != "UPDATE_CODES" {
		log.Printf("UpdateFederationController - Invalid OperationType for federation %s: %s", federationContextId, patchParams.OperationType)
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid OperationType, must be either 'ADD_CODES', 'REMOVE_CODES' or 'UPDATE_CODES'")
		return
	}

	log.Printf("UpdateFederationController - Applying patch to federation: %s", federationContextId)
	err := fmc.federationService.PatchFederation(federationContextId, patchParams)
	if err != nil {
		log.Printf("UpdateFederationController - Error updating federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error updating the Federation object: "+err.Error())
		return
	}
	
	log.Printf("UpdateFederationController - Federation %s updated successfully", federationContextId)

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
	// Check if the a federation exists with the given accessToken
	accessToken := c.Query("accessToken")
	log.Printf("GetFederationContextIdentifierController - Looking up federation by access token")
	
	if !fmc.federationService.ExistsFederationWithAccessToken(accessToken) {
		log.Printf("GetFederationContextIdentifierController - No federation found with provided access token")
		utils.HandleProblem(c, http.StatusBadRequest, "No Federation found with the given accessToken")
		return
	}

	// Retrieve the federation through the accessToken
	log.Printf("GetFederationContextIdentifierController - Retrieving federation by access token")
	federation, err := fmc.federationService.GetFederationByAccessToken(accessToken)
	if err != nil {
		log.Printf("GetFederationContextIdentifierController - Error retrieving federation by access token: %v", err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving federation information: "+err.Error())
		return
	}
	
	log.Printf("GetFederationContextIdentifierController - Found federation with contextId: %s", federation.PartnerOP.FederationContextId)

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
	// Check if the federation exists
	federationContextId := c.Param("federationContextId")
	log.Printf("GetFederationHealthController - Checking health for federation: %s", federationContextId)
	
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		log.Printf("GetFederationHealthController - No federation found with contextId: %s", federationContextId)
		utils.HandleProblem(c, http.StatusBadRequest, "No Federation found with the given federationContextId")
		return
	}

	// Retrieve the federation information
	log.Printf("GetFederationHealthController - Retrieving federation information for: %s", federationContextId)
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("GetFederationHealthController - Error retrieving federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving federation information")
		return
	}
	
	log.Printf("GetFederationHealthController - Successfully retrieved health info for federation: %s", federationContextId)

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
	federationContextId := c.Param("federationContextId")
	log.Printf("RenewFederationController - Starting federation renewal for contextId: %s", federationContextId)
	
	// Check if the federation exists
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		log.Printf("RenewFederationController - No federation found with contextId: %s", federationContextId)
		utils.HandleProblem(c, http.StatusBadRequest, "No Federation found with the given federationContextId")
		return
	}

	// Get the federation object
	log.Printf("RenewFederationController - Retrieving federation object for: %s", federationContextId)
	federation, err := fmc.federationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("RenewFederationController - Error retrieving federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving the Federation object")
		return
	}

	// Renew the federation relationship
	log.Printf("RenewFederationController - Renewing federation: %s", federationContextId)
	federation, err = fmc.federationService.RenewFederation(federation)
	if err != nil {
		log.Printf("RenewFederationController - Error renewing federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error renewing the Federation object")
		return
	}
	
	log.Printf("RenewFederationController - Federation %s renewed successfully. New renewal date: %v, expiry date: %v", federationContextId, federation.PartnerOP.FederationRenewalDate, federation.PartnerOP.FederationExpiryDate)

	response := dto.FederationRenewalResponseData{
		FederationContextId:   federation.PartnerOP.FederationContextId,
		FederationRenewalDate: federation.PartnerOP.FederationRenewalDate,
		FederationExpiryDate:  federation.PartnerOP.FederationExpiryDate,
	}

	c.JSON(http.StatusOK, response)
}
