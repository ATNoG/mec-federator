package ewbi

import (
	"log"
	"log/slog"
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
	zoneService         *services.ZoneService
}

func NewApplicationInstanceLifecycleManagementController(orchestratorService *services.OrchestratorService, artefactService *services.ArtefactService, appInstanceService *services.AppInstanceService, zoneService *services.ZoneService) *ApplicationInstanceLifecycleManagementController {
	return &ApplicationInstanceLifecycleManagementController{
		orchestratorService: orchestratorService,
		artefactService:     artefactService,
		appInstanceService:  appInstanceService,
		zoneService:         zoneService,
	}
}

// @Summary Create an application instance
// @Description Used by origin OP to create an application instance
// @Tags EWBI - ApplicationInstanceLifecycleManagement
func (amc *ApplicationInstanceLifecycleManagementController) CreateAppInstanceController(c *gin.Context) {
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

	// get the vim id of the zone
	vimId, err := amc.zoneService.GetVimId(request.ZoneInfo)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting vim id: "+err.Error())
		return
	}

	// instantiate the appPkg
	appInstanceId, err := amc.orchestratorService.InstantiateAppPkg(artefact.AppPkgId, vimId, request.Config)
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
		ArtefactId:          artefact.Id,
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
func (amc *ApplicationInstanceLifecycleManagementController) DeleteAppInstanceController(c *gin.Context) {
	log.Print("DeleteApplicationInstanceController - Deleting application instance")

	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	// delete the appInstance
	err := amc.orchestratorService.TerminateAppi(appInstanceId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error terminating application instance: "+err.Error())
		return
	}

	// delete the appInstance from the database
	err = amc.appInstanceService.RemoveAppInstance(federationContextId, appInstanceId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error removing application instance: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"appInstanceId": appInstanceId})
}

// @Summary Get details about an application instance
// @Description Used by origin OP to get details about an application instance
// @Tags EWBI - ApplicationInstanceLifecycleManagement
func (amc *ApplicationInstanceLifecycleManagementController) GetAppInstanceDetailsController(c *gin.Context) {
	log.Print("GetAppInstanceDetailsController - Getting application instance details")

	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the appInstance from the orchestrator
	appInstance, err := amc.orchestratorService.GetAppi(appInstanceId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, appInstance)
}

// @Summary Enable application instance KDU
// @Description Used by origin OP to enable an application instance KDU of a certain application instance
// @Tags EWBI - ApplicationInstanceLifecycleManagement
func (amc *ApplicationInstanceLifecycleManagementController) EnableAppInstanceKDUController(c *gin.Context) {
	log.Print("EnableAppInstanceKDUController - Enabling application instance KDU")

	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	// get the rest of the details from the request body
	var request dto.EnableAppInstanceKDURequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	// get the appInstance from the database
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the appInstance from the orchestrator
	orchAppI, err := amc.orchestratorService.GetAppi(appInstance.Id)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the appPkg from the orchestrator database
	appPkg, err := amc.orchestratorService.GetAppPkg(appInstance.AppPkgId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application package: "+err.Error())
		return
	}

	// iterate over the instances and check if the kdu id is in one of the instances
	var nsId string
	instances := orchAppI.Instances[orchAppI.Domain]
	for _, instance := range instances {
		if _, ok := instance.KDUs[request.KDUId]; ok {
			nsId = instance.NSID
			break
		}
	}

	// enable the kdu
	slog.Info("Enabling KDU", "appdId", appPkg.AppdId, "kduId", request.KDUId, "nsId", nsId, "node", request.Node)
	err = amc.orchestratorService.EnableAppInstanceKDU(appPkg.AppdId, request.KDUId, nsId, request.Node)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error enabling application instance KDU: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"appInstance": appInstance})
}

// @Summary Disable application instance KDU
// @Description Used by origin OP to disable an application instance KDU of a certain application instance
// @Tags EWBI - ApplicationInstanceLifecycleManagement
func (amc *ApplicationInstanceLifecycleManagementController) DisableAppInstanceKDUController(c *gin.Context) {
	log.Print("DisableAppInstanceKDUController - Disabling application instance KDU")

	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	// get the rest of the details from the request body
	var request dto.DisableAppInstanceKDURequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	// get the appInstance from the database
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the appInstance from the orchestrator
	orchAppI, err := amc.orchestratorService.GetAppi(appInstance.Id)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the appPkg from the orchestrator database
	appPkg, err := amc.orchestratorService.GetAppPkg(appInstance.AppPkgId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application package: "+err.Error())
		return
	}

	// get the instances in the domain
	instances := orchAppI.Instances[orchAppI.Domain]

	// iterate over the instances and check if the kdu id is in one of the instances
	for _, instance := range instances {
		if _, ok := instance.KDUs[request.KDUId]; ok {
			// disable the kdu
			err = amc.orchestratorService.DisableAppiKDU(appPkg.AppdId, request.KDUId, instance.NSID)
			if err != nil {
				utils.HandleProblem(c, http.StatusInternalServerError, "Error disabling application instance KDU: "+err.Error())
				return
			}

			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"kduId": request.KDUId})
}
