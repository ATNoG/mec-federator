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

/*
 * CreateFederationController
 *	receives an instance of FederationRequestData and creates a new federation relationship with a partner OP
 *	responds with an instance of FederationResponseData
 */
func (fmc *FederationManagementController) CreateFederationController(c *gin.Context) {
	var federationRequestData models.FederationRequestData
	if err := c.ShouldBindJSON(&federationRequestData); err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusBadRequest)
		c.AbortWithStatusJSON(http.StatusBadRequest, problemDetails)
		return
	}

	federation, err := fmc.federationService.CreateFederation(federationRequestData)
	if err != nil {
		problemDetails := utils.NewProblemDetails(http.StatusInternalServerError)
		c.AbortWithStatusJSON(http.StatusInternalServerError, problemDetails)
		return
	}

	federationResponseData := federation.PartnerOP

	c.JSON(http.StatusOK, federationResponseData)
}

/*
 * DeleteFederationController
 *	receives a federationContextId and deletes the federation relationship with the Originating OP
 */
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
