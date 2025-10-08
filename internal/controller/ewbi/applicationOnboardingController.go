package ewbi

/*
 *
 * This file contains the implementation of the Application Onboarding Controller over the E/WBI
 * Exposes Application Onboarding functionalities to manage applications in federations.
 *
 */

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type ApplicationOnboardingController struct {
	federationService   *services.FederationService
	orchestratorService *services.OrchestratorService
	artefactService     *services.ArtefactService
	applicationService  *services.ApplicationService
	zoneService         *services.ZoneService
}

// NewApplicationOnboardingController creates a new instance of the ApplicationOnboardingController
func NewApplicationOnboardingController(federationService *services.FederationService, orchestratorService *services.OrchestratorService, artefactService *services.ArtefactService, applicationService *services.ApplicationService, zoneService *services.ZoneService) *ApplicationOnboardingController {
	return &ApplicationOnboardingController{
		federationService:   federationService,
		orchestratorService: orchestratorService,
		artefactService:     artefactService,
		applicationService:  applicationService,
		zoneService:         zoneService,
	}
}

// @Summary Onboard Application
// @Description Onboards an application to the orchestrator and creates an application object in the federation
// @Tags EWBI - ApplicationOnboarding
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param request body dto.OnboardApplicationRequest true "Onboard Application Request"
// @Success 201 {object} map[string]string "message: Application onboarded successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request or validation error"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/applications [post]
func (aoc *ApplicationOnboardingController) OnboardApplicationController(c *gin.Context) {
	federationContextId := c.Param("federationContextId")
	log.Printf("OnboardApplicationController - Starting application onboarding for federation: %s", federationContextId)

	// get and bind the request body
	log.Printf("OnboardApplicationController - Binding request body for federation: %s", federationContextId)
	var request dto.OnboardApplicationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("OnboardApplicationController - Error binding request body for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	// get the artefact from the database
	log.Printf("OnboardApplicationController - Retrieving artefact for federation: %s, artefactId: %s", federationContextId, request.AppId)
	artefact, err := aoc.artefactService.GetArtefact(federationContextId, request.ArtefactId)
	if err != nil {
		log.Printf("OnboardApplicationController - Error getting artefact from database for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting artefact: "+err.Error())
		return
	}

	descriptorData, err := utils.GetDescriptorData(*artefact.ArtefactFile)
	if err != nil {
		log.Printf("OnboardApplicationController - Error getting descriptor data for federation %s, artefactId %s: %v", federationContextId, artefact.Id, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting descriptor data: "+err.Error())
		return
	}

	// validate the descriptor data
	log.Printf("OnboardApplicationController - Validating descriptor data for federation: %s, artefactId: %s", federationContextId, artefact.Id)
	appPkg, err := utils.ValidateDescriptorData(descriptorData)
	if err != nil {
		log.Printf("OnboardApplicationController - Error validating descriptor data for federation %s, artefactId %s: %v", federationContextId, artefact.Id, err)
		utils.HandleProblem(c, http.StatusBadRequest, "Error validating descriptor data: "+err.Error())
		return
	}
	appPkg.AppD = *artefact.ArtefactFile

	// onboard the artefact into the orchestrator
	log.Printf("OnboardApplicationController - Onboarding artefact to orchestrator for federation: %s, artefactId: %s", federationContextId, artefact.Id)
	appPkgId, err := aoc.orchestratorService.OnboardAppPkg(appPkg)
	if err != nil {
		log.Printf("OnboardApplicationController - Error onboarding artefact to orchestrator for federation %s, artefactId %s: %v", federationContextId, artefact.Id, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error onboarding artefact onto orchestrator: "+err.Error())
		return
	}

	// update the artefact object
	log.Printf("OnboardApplicationController - Updating artefact object for federation: %s, artefactId: %s", federationContextId, artefact.Id)
	artefact.AppPkgId = appPkgId
	err = aoc.artefactService.UpdateArtefact(artefact)
	if err != nil {
		log.Printf("OnboardApplicationController - Error updating artefact object for federation %s, artefactId %s: %v", federationContextId, artefact.Id, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error updating artefact object: "+err.Error())
		return
	}

	log.Printf("OnboardApplicationController - Artefact object updated for federation: %s, artefactId: %s", federationContextId, artefact.Id)
	log.Printf("OnboardApplicationController - Application onboarded successfully for federation: %s, artefactId: %s", federationContextId, artefact.Id)

	// create the application object
	log.Printf("OnboardApplicationController - Creating application object for federation: %s, appId: %s", federationContextId, request.AppId)
	application := models.Application{
		Id:                  request.AppId,
		Name:                request.AppName,
		FederationContextId: federationContextId,
		ArtefactId:          artefact.Id,
		AppPkgId:            appPkgId,
	}

	err = aoc.applicationService.CreateApplication(application)
	if err != nil {
		log.Printf("OnboardApplicationController - Error creating application object for federation %s, appId %s: %v", federationContextId, request.AppId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error creating application object: "+err.Error())
		return
	}

	log.Printf("OnboardApplicationController - Application object created for federation: %s, appId: %s", federationContextId, request.AppId)
	log.Printf("OnboardApplicationController - Application onboarded successfully for federation: %s, appId: %s", federationContextId, request.AppId)

	c.JSON(http.StatusCreated, gin.H{"message": "Application onboarded successfully"})
}

// @Summary Update Application
// @Description Updates an application object in the federation
// @Tags EWBI - ApplicationOnboarding
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param appId path string true "Application ID"
// @Success 200 {object} map[string]string "message: Application updated successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/applications/{appId} [patch]
func (aoc *ApplicationOnboardingController) UpdateApplicationController(c *gin.Context) {
}

// @Summary Remove Application
// @Description Removes an application from the orchestrator and deletes the application object from the federation
// @Tags EWBI - ApplicationOnboarding
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param appId path string true "Application ID"
// @Success 200 {object} map[string]string "message: Application removed successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request or application has running instances"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/applications/{appId} [delete]
func (aoc *ApplicationOnboardingController) RemoveApplicationController(c *gin.Context) {
	federationContextId := c.Param("federationContextId")

	appId := c.Param("appId")

	log.Printf("RemoveApplicationController - Removing application for federation: %s, appId: %s", federationContextId, appId)

	// get the application object
	log.Printf("RemoveApplicationController - Retrieving application object for federation: %s, appId: %s", federationContextId, appId)
	application, err := aoc.applicationService.GetApplication(federationContextId, appId)
	if err != nil {
		log.Printf("RemoveApplicationController - Error retrieving application object for federation %s, appId %s: %v", federationContextId, appId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error retrieving application object: "+err.Error())
		return
	}

	// TODO: check if there are no app instances of this application

	// remove the app package from the orchestrator
	log.Printf("RemoveApplicationController - Removing app package from orchestrator for federation: %s, appId: %s", federationContextId, appId)
	err = aoc.orchestratorService.RemoveAppPkg(application.AppPkgId)
	if err != nil {
		log.Printf("RemoveApplicationController - Error removing app package from orchestrator for federation %s, appId %s: %v", federationContextId, appId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error removing app package from orchestrator: "+err.Error())
		return
	}

	// remove the application object
	log.Printf("RemoveApplicationController - Removing application object for federation: %s, appId: %s", federationContextId, appId)
	err = aoc.applicationService.DeleteApplication(federationContextId, appId)
	if err != nil {
		log.Printf("RemoveApplicationController - Error removing application object for federation %s, appId %s: %v", federationContextId, appId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error removing application object: "+err.Error())
		return
	}

	log.Printf("RemoveApplicationController - Application object removed for federation: %s, appId: %s", federationContextId, appId)
	log.Printf("RemoveApplicationController - Application removed successfully for federation: %s, appId: %s", federationContextId, appId)

	c.JSON(http.StatusOK, gin.H{"message": "Application removed successfully"})
}

// @Summary View Application
// @Description Retrieves application details by application ID within a federation
// @Tags EWBI - ApplicationOnboarding
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param appId path string true "Application ID"
// @Success 200 {object} models.Application
// @Failure 400 {object} models.ProblemDetails "Invalid request or application not found"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/applications/{appId} [get]
func (aoc *ApplicationOnboardingController) ViewApplicationController(c *gin.Context) {
}
