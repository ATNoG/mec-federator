package ewbi

/*
 *
 * This file contains the implementation of the Federation Management Controller
 *
 */

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
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

// @Summary Create Federation Relationship
// @Description Creates a new federation with another federator with the provided data
// @Tags FederationManagement
// @Accept json
// @Produce json
// @Param federationRequestData body models.FederationRequestData true "Federation Request Data"
// @Success 200 {object} models.FederationResponseData
// @Failure 400 {object} models.ProblemDetails
// @Failure 500 {object} models.ProblemDetails
// @Router /federation/v1/partner [post]
func (fmc *FederationManagementController) CreateFederationController(c *gin.Context) {
	var federationRequestData models.FederationRequestData
	if err := c.ShouldBindJSON(&federationRequestData); err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		c.AbortWithStatusJSON(http.StatusBadRequest, problemDetails)
		return
	}

	// need to update this endpoint to relate access tokens to the federation

	federation, err := fmc.federationService.CreateFederation(federationRequestData)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		c.AbortWithStatusJSON(http.StatusInternalServerError, problemDetails)
		return
	}

	c.JSON(http.StatusOK, federation.PartnerOP)
}

// @Summary Get Federation Meta Information
// @Description Retrieves metadata information about a federation based on the federationContextId
// @Tags FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} models.FederationMetaInfo
// @Failure 400 {object} models.ProblemDetails "Invalid Federation Context ID"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /federation/v1/{federationContextId}/partner [get]
func (fmc *FederationManagementController) GetFederationMetaInfoController(c *gin.Context) {
	federationContextId := c.Param("federationContextId")

	federation, err := fmc.federationService.GetFederationFromContextId(federationContextId)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "No Federation found with the given federationContextId"
		c.AbortWithStatusJSON(http.StatusInternalServerError, problemDetails)
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

// @Summary Remove Federation Relationship
// @Description Removes a federation relationship by its federationContextId
// @Tags FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {object} map[string]string "status: Federation removed successfully"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /federation/v1/{federationContextId}/partner [delete]
func (fmc *FederationManagementController) RemoveFederationRelationshipController(c *gin.Context) {
	// Remove the federation object from the database
	federationContextId := c.Param("federationContextId")

	err := fmc.federationService.DeleteFederation(federationContextId)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		c.AbortWithStatusJSON(http.StatusInternalServerError, problemDetails)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Federation removed successfully"})
}

// @Summary Update a Federation
// @Description Updates a federation object with the given federationContextId and patch parameters
// @Tags FederationManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param patchParams body models.FederationPatchParams true "Patch Parameters"
// @Success 200 {object} map[string]string "status: Federation updated successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request or federation not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /federation/v1/{federationContextId}/partner [patch]
func (fmc *FederationManagementController) UpdateFederationController(c *gin.Context) {
	// Check if the request body is valid and within expected data type
	var patchParams models.FederationPatchParams
	if err := c.ShouldBindJSON(&patchParams); err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "Invalid request body, missing parameters or wrong data type"
		c.AbortWithStatusJSON(http.StatusBadRequest, problemDetails)
		return
	}

	// Get the federationContextId from the path
	federationContextId := c.Param("federationContextId")
	// Check if the federation exists
	if !fmc.federationService.ExistsFederationWithContextId(federationContextId) {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "No Federation found with the given federationContextId"
		c.AbortWithStatusJSON(http.StatusBadRequest, problemDetails)
		return
	}

	// Check validity of patchParams
	if patchParams.ObjectType != "MOBILE_NETWORK_CODES" && patchParams.ObjectType != "FIXED_NETWORK_CODES" {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "Invalid ObjectType, must be either 'MOBILE_NETWORK_CODES' or 'FIXED_NETWORK_CODES'"
		c.AbortWithStatusJSON(http.StatusBadRequest, problemDetails)
		return
	}

	if patchParams.OperationType != "ADD_CODES" && patchParams.OperationType != "REMOVE_CODES" && patchParams.OperationType != "UPDATE_CODES" {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		problemDetails.Detail = "Invalid OperationType, must be either 'ADD_CODES', 'REMOVE_CODES' or 'UPDATE_CODES'"
		c.AbortWithStatusJSON(http.StatusBadRequest, problemDetails)
		return
	}

	// Update the federation object
	err := fmc.federationService.PatchFederation(federationContextId, patchParams)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		problemDetails.Detail = "Error updating the Federation object"
		c.AbortWithStatusJSON(http.StatusInternalServerError, problemDetails)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Federation updated successfully"})
}
