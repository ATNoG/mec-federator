package ewbi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type ApplicationInstanceLifecycleManagementController struct {
	orchestratorService *services.OrchestratorService
	artefactService     *services.ArtefactService
	appInstanceService  *services.AppInstanceService
}

func NewApplicationInstanceLifecycleManagementController(orchestratorService *services.OrchestratorService, artefactService *services.ArtefactService, appInstanceService *services.AppInstanceService) *ApplicationInstanceLifecycleManagementController {
	return &ApplicationInstanceLifecycleManagementController{
		orchestratorService: orchestratorService,
		artefactService:     artefactService,
		appInstanceService:  appInstanceService,
	}
}

// @Summary Create an application instance
// @Description Used by origin OP to create an application instance
// @Tags EWBI - ApplicationInstanceLifecycleManagement
func (amc *ApplicationInstanceLifecycleManagementController) CreateApplicationInstanceController(c *gin.Context) {
	log.Print("CreateApplicationInstanceController - Creating application instance")

	// get and bind the request body
	var request dto.InstantiateApplicationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	// get the artefact from the database
	federationContextId := c.Param("federationContextId")
	artefact, err := amc.artefactService.GetArtefact(federationContextId, request.AppId)
	if err != nil {
		utils.HandleProblem(c, http.StatusNotFound, "Artefact not found: "+err.Error())
		return
	}

	// instantiate the appPkg
	appInstanceId, err := amc.orchestratorService.InstantiateAppPkg(artefact.AppPkgId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error instantiating application instance: "+err.Error())
		return
	}

	// create the appInstance
	appInstance := models.AppInstance{
		Id:                  appInstanceId,
		FederationContextId: federationContextId,
		Name:                request.AppId,
		Description:         request.AppVersion,
		AppPkgId:            artefact.AppPkgId,
	}

	// save the appInstance to the database
	err = amc.appInstanceService.RegisterAppInstance(appInstance)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error creating application instance: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"appInstanceId": appInstanceId})

}

// @Summary Delete an application instance
// @Description Used by origin OP to delete an application instance
// @Tags EWBI - ApplicationInstanceLifecycleManagement
func (amc *ApplicationInstanceLifecycleManagementController) DeleteApplicationInstanceController(c *gin.Context) {
	log.Print("DeleteApplicationInstanceController - Deleting application instance")

	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// delete the appInstance
	err := amc.orchestratorService.TerminateAppPkg(appInstanceId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error terminating application instance: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"appInstanceId": appInstanceId})
}

func (amc *ApplicationInstanceLifecycleManagementController) GetAppInstanceController(c *gin.Context) {
	log.Print("GetAppInstanceController - Getting application instance")

	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the appInstance from the orchestrator
	appInstance, err := amc.orchestratorService.GetAppInstance(appInstanceId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, appInstance)
}
