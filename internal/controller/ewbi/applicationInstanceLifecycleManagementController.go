package ewbi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Description Creates a new application instance by instantiating an artefact in the specified zone. The operation includes artefact validation, orchestrator instantiation, and database registration.
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Param federationContextId path string true "Federation Context ID" format(uuid)
// @Accept json
// @Produce json
// @Param request body dto.InstantiateApplicationRequest true "Application instantiation request"
// @Success 201 {object} map[string]string "Application instance created successfully with appInstanceId and nsId"
// @Failure 400 {object} models.ProblemDetails "Bad request - invalid request body or validation errors"
// @Failure 404 {object} models.ProblemDetails "Artefact not found in the specified federation context"
// @Failure 500 {object} models.ProblemDetails "Internal server error - VIM ID retrieval, orchestrator instantiation, or database errors"
// @Router /ewbi/{federationContextId}/app_instances [post]
func (amc *ApplicationInstanceLifecycleManagementController) CreateAppInstanceController(c *gin.Context) {
	federationContextId := c.Param("federationContextId")
	log.Printf("CreateAppInstanceController - Starting application instance creation for federation: %s", federationContextId)

	// get and bind the request body
	log.Printf("CreateAppInstanceController - Binding request body for federation: %s", federationContextId)
	var request dto.InstantiateApplicationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("CreateAppInstanceController - Error binding request body for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	// get the artefact from the database
	log.Printf("CreateAppInstanceController - Retrieving artefact for federation: %s, appId: %s", federationContextId, request.AppId)
	artefact, err := amc.artefactService.GetArtefact(federationContextId, request.AppId)
	if err != nil {
		log.Printf("CreateAppInstanceController - Artefact not found for federation %s, appId %s: %v", federationContextId, request.AppId, err)
		utils.HandleProblem(c, http.StatusNotFound, "Artefact not found: "+err.Error())
		return
	}

	// get the vim id of the zone
	log.Printf("CreateAppInstanceController - Getting VIM ID for federation: %s, appId: %s, zone: %v", federationContextId, request.AppId, request.ZoneInfo)
	vimId, err := amc.zoneService.GetVimId(request.ZoneInfo)
	if err != nil {
		log.Printf("CreateAppInstanceController - Error getting VIM ID for federation %s, appId %s: %v", federationContextId, request.AppId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting vim id: "+err.Error())
		return
	}

	log.Printf("CreateAppInstanceController - Retrieved VIM ID for federation: %s, appId: %s, vimId: %s", federationContextId, request.AppId, vimId)

	// instantiate the appPkg
	log.Printf("CreateAppInstanceController - Instantiating app package for federation: %s, appId: %s, appPkgId: %s, vimId: %s", federationContextId, request.AppId, artefact.AppPkgId, vimId)
	appiId, err := amc.orchestratorService.InstantiateAppPkg(artefact.AppPkgId, vimId, request.Config)
	if err != nil {
		log.Printf("CreateAppInstanceController - Error instantiating app package for federation %s, appId %s: %v", federationContextId, request.AppId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error instantiating application instance: "+err.Error())
		return
	}

	log.Printf("CreateAppInstanceController - App package instantiated successfully for federation: %s, appId: %s, appiId: %s", federationContextId, request.AppId, appiId)

	// create the appInstance
	appInstanceId := uuid.New().String()
	log.Printf("CreateAppInstanceController - Creating app instance object for federation: %s, appId: %s, appInstanceId: %s", federationContextId, request.AppId, appInstanceId)
	appInstance := models.AppInstance{
		Id:                  appInstanceId,
		FederationContextId: federationContextId,
		Name:                "federated-instance",
		Description:         "local",
		ArtefactId:          artefact.Id,
		AppiId:              appiId,
		AppPkgId:            artefact.AppPkgId,
	}

	// save the appInstance to the database
	log.Printf("CreateAppInstanceController - Saving app instance to database for federation: %s, appInstanceId: %s", federationContextId, appInstance.Id)
	err = amc.appInstanceService.RegisterAppInstance(appInstance)
	if err != nil {
		log.Printf("CreateAppInstanceController - Error saving app instance to database for federation %s, appInstanceId %s: %v", federationContextId, appInstance.Id, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error creating application instance: "+err.Error())
		return
	}

	// get the appi from the orchestrator
	log.Printf("CreateAppInstanceController - Retrieving orchestrator app instance for federation: %s, appInstanceId: %s, appiId: %s", federationContextId, appInstance.Id, appiId)
	orchAppI, err := amc.orchestratorService.GetAppi(appiId)
	if err != nil {
		log.Printf("CreateAppInstanceController - Error getting orchestrator app instance for federation %s, appInstanceId %s: %v", federationContextId, appInstance.Id, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the zone from the vim id
	zone, err := amc.zoneService.GetZoneFromVimId(vimId)
	if err != nil {
		log.Printf("CreateAppInstanceController - Error getting zone from vim id for federation %s, appInstanceId %s: %v", federationContextId, appInstance.Id, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting zone: "+err.Error())
		return
	}

	nsId := orchAppI.Instances[orchAppI.Domain][zone.ZoneId].NSID
	log.Printf("CreateAppInstanceController - Retrieved NS ID for federation: %s, appInstanceId: %s, nsId: %s", federationContextId, appInstance.Id, nsId)
	log.Printf("CreateAppInstanceController - App instance created successfully for federation: %s, appInstanceId: %s, nsId: %s", federationContextId, appInstance.Id, nsId)
	c.JSON(http.StatusCreated, gin.H{"appInstanceId": appInstance.Id, "nsId": nsId})
}

// @Summary Delete an application instance
// @Description Terminates and removes an application instance from both the orchestrator and database. This operation cleans up all associated resources.
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Param federationContextId path string true "Federation Context ID" format(uuid)
// @Param appInstanceId path string true "Application Instance ID" format(uuid)
// @Produce json
// @Success 200 {object} map[string]string "Application instance deleted successfully with appInstanceId"
// @Failure 404 {object} models.ProblemDetails "Application instance not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error - orchestrator termination or database removal errors"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId} [delete]
func (amc *ApplicationInstanceLifecycleManagementController) DeleteAppInstanceController(c *gin.Context) {
	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	log.Printf("DeleteAppInstanceController - Starting app instance deletion for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)

	// get the appInstance from the database
	log.Printf("DeleteAppInstanceController - Getting app instance from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("DeleteAppInstanceController - Error getting app instance from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// delete the appInstance
	log.Printf("DeleteAppInstanceController - Terminating app instance in orchestrator for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	err = amc.orchestratorService.TerminateAppi(appInstance.AppiId)
	if err != nil {
		log.Printf("DeleteAppInstanceController - Error terminating app instance in orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error terminating application instance: "+err.Error())
		return
	}

	// delete the appInstance from the database
	log.Printf("DeleteAppInstanceController - Removing app instance from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	err = amc.appInstanceService.RemoveAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("DeleteAppInstanceController - Error removing app instance from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error removing application instance: "+err.Error())
		return
	}

	log.Printf("DeleteAppInstanceController - App instance deleted successfully for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	c.JSON(http.StatusOK, gin.H{"appInstanceId": appInstanceId})
}

// @Summary Get application instance details
// @Description Retrieves comprehensive details about a specific application instance from the orchestrator, including deployment status and configuration.
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Param federationContextId path string true "Federation Context ID" format(uuid)
// @Param appInstanceId path string true "Application Instance ID" format(uuid)
// @Produce json
// @Success 200 {object} dto.OrchAppI "Application instance details from orchestrator"
// @Failure 404 {object} models.ProblemDetails "Application instance not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error - orchestrator access failure"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId} [get]
func (amc *ApplicationInstanceLifecycleManagementController) GetAppInstanceDetailsController(c *gin.Context) {
	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")
	federationContextId := c.Param("federationContextId")

	log.Printf("GetAppInstanceDetailsController - Getting app instance details for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)

	// get the appInstance from the orchestrator
	log.Printf("GetAppInstanceDetailsController - Retrieving app instance from orchestrator for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	appInstance, err := amc.orchestratorService.GetAppi(appInstanceId)
	if err != nil {
		log.Printf("GetAppInstanceDetailsController - Error getting app instance from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	log.Printf("GetAppInstanceDetailsController - Successfully retrieved app instance details for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	c.JSON(http.StatusOK, appInstance)
}

// @Summary Enable application instance KDU
// @Description Enables a specific Kubernetes Deployment Unit (KDU) within an application instance. This operation activates the KDU on the specified node.
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Param federationContextId path string true "Federation Context ID" format(uuid)
// @Param appInstanceId path string true "Application Instance ID" format(uuid)
// @Accept json
// @Produce json
// @Param request body dto.EnableAppInstanceKDURequest true "KDU enablement request with KDU ID and target node"
// @Success 200 {object} map[string]interface{} "KDU enabled successfully with application instance details"
// @Failure 400 {object} models.ProblemDetails "Bad request - invalid request body or missing KDU"
// @Failure 404 {object} models.ProblemDetails "Application instance or KDU not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error - database access, orchestrator operations, or KDU enablement failure"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId}/kdu/enable [post]
func (amc *ApplicationInstanceLifecycleManagementController) EnableAppInstanceKDUController(c *gin.Context) {
	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	log.Printf("EnableAppInstanceKDUController - Starting KDU enablement for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)

	// get the rest of the details from the request body
	log.Printf("EnableAppInstanceKDUController - Binding request body for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	var request dto.EnableAppInstanceKDURequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("EnableAppInstanceKDUController - Error binding request body for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("EnableAppInstanceKDUController - Request details for federation: %s, appInstanceId: %s, kduId: %s, node: %s", federationContextId, appInstanceId, request.KDUId, request.Node)

	// get the appInstance from the database
	log.Printf("EnableAppInstanceKDUController - Getting app instance from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("EnableAppInstanceKDUController - Error getting app instance from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the appInstance from the orchestrator
	log.Printf("EnableAppInstanceKDUController - Getting app instance from orchestrator for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	orchAppI, err := amc.orchestratorService.GetAppi(appInstance.Id)
	if err != nil {
		log.Printf("EnableAppInstanceKDUController - Error getting app instance from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the appPkg from the orchestrator database
	log.Printf("EnableAppInstanceKDUController - Getting app package from orchestrator for federation: %s, appInstanceId: %s, appPkgId: %s", federationContextId, appInstanceId, appInstance.AppPkgId)
	appPkg, err := amc.orchestratorService.GetAppPkg(appInstance.AppPkgId)
	if err != nil {
		log.Printf("EnableAppInstanceKDUController - Error getting app package from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application package: "+err.Error())
		return
	}

	// iterate over the instances and check if the kdu id is in one of the instances
	log.Printf("EnableAppInstanceKDUController - Searching for KDU in instances for federation: %s, appInstanceId: %s, kduId: %s", federationContextId, appInstanceId, request.KDUId)
	var nsId string
	instances := orchAppI.Instances[orchAppI.Domain]
	for _, instance := range instances {
		if _, ok := instance.KDUs[request.KDUId]; ok {
			nsId = instance.NSID
			log.Printf("EnableAppInstanceKDUController - Found KDU in instance for federation: %s, appInstanceId: %s, kduId: %s, nsId: %s", federationContextId, appInstanceId, request.KDUId, nsId)
			break
		}
	}

	// enable the kdu
	log.Printf("EnableAppInstanceKDUController - Enabling KDU for federation: %s, appInstanceId: %s, kduId: %s, nsId: %s, node: %s, appdId: %s", federationContextId, appInstanceId, request.KDUId, nsId, request.Node, appPkg.AppdId)
	err = amc.orchestratorService.EnableAppInstanceKDU(appPkg.AppdId, request.KDUId, nsId, request.Node)
	if err != nil {
		log.Printf("EnableAppInstanceKDUController - Error enabling KDU for federation %s, appInstanceId %s, kduId %s: %v", federationContextId, appInstanceId, request.KDUId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error enabling application instance KDU: "+err.Error())
		return
	}

	log.Printf("EnableAppInstanceKDUController - KDU enabled successfully for federation: %s, appInstanceId: %s, kduId: %s", federationContextId, appInstanceId, request.KDUId)
	c.JSON(http.StatusOK, gin.H{"appInstance": appInstance})
}

// @Summary Disable application instance KDU
// @Description Disables a specific Kubernetes Deployment Unit (KDU) within an application instance. This operation deactivates the KDU and stops its execution.
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Param federationContextId path string true "Federation Context ID" format(uuid)
// @Param appInstanceId path string true "Application Instance ID" format(uuid)
// @Accept json
// @Produce json
// @Param request body dto.DisableAppInstanceKDURequest true "KDU disablement request with KDU ID"
// @Success 200 {object} map[string]string "KDU disabled successfully with KDU ID"
// @Failure 400 {object} models.ProblemDetails "Bad request - invalid request body or missing KDU"
// @Failure 404 {object} models.ProblemDetails "Application instance or KDU not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error - database access, orchestrator operations, or KDU disablement failure"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId}/kdu/disable [post]
func (amc *ApplicationInstanceLifecycleManagementController) DisableAppInstanceKDUController(c *gin.Context) {
	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	log.Printf("DisableAppInstanceKDUController - Starting KDU disablement for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)

	// get the rest of the details from the request body
	log.Printf("DisableAppInstanceKDUController - Binding request body for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	var request dto.DisableAppInstanceKDURequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("DisableAppInstanceKDUController - Error binding request body for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("DisableAppInstanceKDUController - Request details for federation: %s, appInstanceId: %s, kduId: %s", federationContextId, appInstanceId, request.KDUId)

	// get the appInstance from the database
	log.Printf("DisableAppInstanceKDUController - Getting app instance from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("DisableAppInstanceKDUController - Error getting app instance from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the appInstance from the orchestrator
	log.Printf("DisableAppInstanceKDUController - Getting app instance from orchestrator for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	orchAppI, err := amc.orchestratorService.GetAppi(appInstance.AppiId)
	if err != nil {
		log.Printf("DisableAppInstanceKDUController - Error getting app instance from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the appPkg from the orchestrator database
	log.Printf("DisableAppInstanceKDUController - Getting app package from orchestrator for federation: %s, appInstanceId: %s, appPkgId: %s", federationContextId, appInstanceId, appInstance.AppPkgId)
	appPkg, err := amc.orchestratorService.GetAppPkg(appInstance.AppPkgId)
	if err != nil {
		log.Printf("DisableAppInstanceKDUController - Error getting app package from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application package: "+err.Error())
		return
	}

	// get the instances in the domain
	log.Printf("DisableAppInstanceKDUController - Searching for KDU in instances for federation: %s, appInstanceId: %s, kduId: %s", federationContextId, appInstanceId, request.KDUId)
	instances := orchAppI.Instances[orchAppI.Domain]

	// iterate over the instances and check if the kdu id is in one of the instances
	for _, instance := range instances {
		if _, ok := instance.KDUs[request.KDUId]; ok {
			log.Printf("DisableAppInstanceKDUController - Found KDU in instance, disabling for federation: %s, appInstanceId: %s, kduId: %s, nsId: %s, appdId: %s", federationContextId, appInstanceId, request.KDUId, instance.NSID, appPkg.AppdId)
			// disable the kdu
			err = amc.orchestratorService.DisableAppiKDU(appPkg.AppdId, request.KDUId, instance.NSID)
			if err != nil {
				log.Printf("DisableAppInstanceKDUController - Error disabling KDU for federation %s, appInstanceId %s, kduId %s: %v", federationContextId, appInstanceId, request.KDUId, err)
				utils.HandleProblem(c, http.StatusInternalServerError, "Error disabling application instance KDU: "+err.Error())
				return
			}

			log.Printf("DisableAppInstanceKDUController - KDU disabled successfully for federation: %s, appInstanceId: %s, kduId: %s", federationContextId, appInstanceId, request.KDUId)
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"kduId": request.KDUId})
}
