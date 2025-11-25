package ewbi

/*
 *
 * This file contains the implementation of the Application Instance Lifecycle Management Controller over the E/WBI
 * Exposes Application Instance Lifecycle Management functionalities to manage app instances in federations.
 *
 */

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type ApplicationInstanceLifecycleManagementController struct {
	federationService   *services.FederationService
	orchestratorService *services.OrchestratorService
	applicationService  *services.ApplicationService
	appInstanceService  *services.AppInstanceService
	zoneService         *services.ZoneService
}

// NewApplicationInstanceLifecycleManagementController creates a new instance of the ApplicationInstanceLifecycleManagementController
func NewApplicationInstanceLifecycleManagementController(federationService *services.FederationService, orchestratorService *services.OrchestratorService, applicationService *services.ApplicationService, appInstanceService *services.AppInstanceService, zoneService *services.ZoneService) *ApplicationInstanceLifecycleManagementController {
	return &ApplicationInstanceLifecycleManagementController{
		federationService:   federationService,
		orchestratorService: orchestratorService,
		applicationService:  applicationService,
		appInstanceService:  appInstanceService,
		zoneService:         zoneService,
	}
}

// @Summary Create Application Instance
// @Description Creates a new application instance in the specified zone and registers it in the orchestrator
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param request body dto.InstantiateApplicationRequest true "Application instantiation request"
// @Success 201 {object} map[string]string "appInstanceId: id, nsId: id, vnfId: id"
// @Failure 400 {object} models.ProblemDetails "Invalid request body"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/app_instances [post]
func (amc *ApplicationInstanceLifecycleManagementController) CreateAppInstanceController(c *gin.Context) {
	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-instantiate-appi-init",
		Message: "",
		Value:   nil,
	})

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

	// get the federation from the database
	log.Printf("CreateAppInstanceController - Getting federation from database for federation: %s", federationContextId)
	federation, err := amc.federationService.GetFederation(federationContextId)
	if err != nil {
		log.Printf("CreateAppInstanceController - Error getting federation from database for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting federation: "+err.Error())
		return
	}

	// get the application from the database
	log.Printf("CreateAppInstanceController - Getting application from database for federation: %s, appId: %s", federationContextId, request.AppId)
	application, err := amc.applicationService.GetApplication(federationContextId, request.AppId)
	if err != nil {
		log.Printf("CreateAppInstanceController - Error getting application from database for federation %s, appId %s: %v", federationContextId, request.AppId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application: "+err.Error())
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
	log.Printf("CreateAppInstanceController - Instantiating app package for federation: %s, appId: %s, appPkgId: %s, vimId: %s", federationContextId, request.AppId, application.AppPkgId, vimId)
	appiId, err := amc.orchestratorService.InstantiateAppPkg(application.AppPkgId, vimId, request.Config, federation.OriginOP.OrigOPFederationId)
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
		AppiId:              appiId,
		AppId:               application.Id,
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

	nsId := orchAppI.Instances[config.AppConfig.OperatorId][zone.ZoneId].NSID
	vnfId := orchAppI.Instances[config.AppConfig.OperatorId][zone.ZoneId].VNFID
	log.Printf("CreateAppInstanceController - App instance created successfully for federation: %s, appInstanceId: %s, nsId: %s, vnfId: %s", federationContextId, appInstance.Id, nsId, vnfId)

	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-instantiate-appi-done",
		Message: "",
		Value:   nil,
	})

	c.JSON(http.StatusCreated, gin.H{"appInstanceId": appInstance.Id, "nsId": nsId, "vnfId": vnfId})
}

// @Summary Delete Application Instance
// @Description Terminates and removes an application instance from the orchestrator and database
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param appInstanceId path string true "Application Instance ID"
// @Success 200 {object} map[string]string "appInstanceId: id"
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId} [delete]
func (amc *ApplicationInstanceLifecycleManagementController) DeleteAppInstanceController(c *gin.Context) {
	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-remove-appi-init",
		Message: "",
		Value:   nil,
	})

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

	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-remove-appi-done",
		Message: "",
		Value:   nil,
	})

	c.JSON(http.StatusOK, gin.H{"appInstanceId": appInstanceId})
}

// @Summary Get Application Instance Details
// @Description Retrieves detailed information about an application instance from the orchestrator
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param appInstanceId path string true "Application Instance ID"
// @Success 200 {object} dto.OrchAppI
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId} [get]
func (amc *ApplicationInstanceLifecycleManagementController) GetAppInstanceDetailsController(c *gin.Context) {
	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	log.Printf("GetAppInstanceDetailsController - Getting app instance details for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)

	// get the appInstance from the orchestrator
	log.Printf("GetAppInstanceDetailsController - Retrieving app instance from orchestrator for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("GetAppInstanceDetailsController - Error getting app instance from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	log.Printf("GetAppInstanceDetailsController - Successfully retrieved app instance details for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	c.JSON(http.StatusOK, appInstance)
}

// @Summary Enable Application Instance KDU
// @Description Enables a Kubernetes Deployment Unit (KDU) within an application instance on the specified node
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param appInstanceId path string true "Application Instance ID"
// @Param request body dto.EnableAppInstanceKDURequest true "KDU enablement request"
// @Success 200 {object} map[string]string "appInstance: object, kduId: id, nsId: id"
// @Failure 400 {object} models.ProblemDetails "Invalid request body"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId}/kdu/enable [post]
func (amc *ApplicationInstanceLifecycleManagementController) EnableAppInstanceKDUController(c *gin.Context) {
	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-enable-kdu-init",
		Message: "",
		Value:   nil,
	})

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

	log.Printf("EnableAppInstanceKDUController - Request details for federation: %s, appInstanceId: %s, kduId: %s, node: %s", federationContextId, appInstanceId, request.KduId, request.Node)

	// get the appInstance from the database
	log.Printf("EnableAppInstanceKDUController - Getting app instance from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("EnableAppInstanceKDUController - Error getting app instance from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the application from the database
	log.Printf("EnableAppInstanceKDUController - Getting application from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	application, err := amc.applicationService.GetApplication(federationContextId, appInstance.AppId)
	if err != nil {
		log.Printf("EnableAppInstanceKDUController - Error getting application from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application: "+err.Error())
		return
	}

	// get the appPkg from the orchestrator database
	log.Printf("EnableAppInstanceKDUController - Getting app package from orchestrator for federation: %s, appInstanceId: %s, appPkgId: %s", federationContextId, appInstanceId, application.AppPkgId)
	appPkg, err := amc.orchestratorService.GetAppPkg(application.AppPkgId)
	if err != nil {
		log.Printf("EnableAppInstanceKDUController - Error getting app package from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application package: "+err.Error())
		return
	}

	// enable the kdu
	log.Printf("EnableAppInstanceKDUController - Enabling KDU for federation: %s, appInstanceId: %s, kduId: %s, nsId: %s, node: %s, appdId: %s", federationContextId, appInstanceId, request.KduId, request.NsId, request.Node, appPkg.AppdId)
	err = amc.orchestratorService.EnableAppInstanceKDU(appPkg.AppdId, request.KduId, request.NsId, request.Node)
	if err != nil {
		log.Printf("EnableAppInstanceKDUController - Error enabling KDU for federation %s, appInstanceId %s, kduId %s: %v", federationContextId, appInstanceId, request.KduId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error enabling application instance KDU: "+err.Error())
		return
	}

	log.Printf("EnableAppInstanceKDUController - KDU enabled successfully for federation: %s, appInstanceId: %s, kduId: %s", federationContextId, appInstanceId, request.KduId)

	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-enable-kdu-done",
		Message: "",
		Value:   nil,
	})

	c.JSON(http.StatusOK, gin.H{"appInstance": appInstance, "kduId": request.KduId, "nsId": request.NsId})
}

// @Summary Disable Application Instance KDU
// @Description Disables a Kubernetes Deployment Unit (KDU) within an application instance
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param appInstanceId path string true "Application Instance ID"
// @Param request body dto.DisableAppInstanceKDURequest true "KDU disablement request"
// @Success 200 {object} map[string]string "appInstance: object, kduId: id, nsId: id"
// @Failure 400 {object} models.ProblemDetails "Invalid request body"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId}/kdu/disable [post]
func (amc *ApplicationInstanceLifecycleManagementController) DisableAppInstanceKDUController(c *gin.Context) {
	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-disable-kdu-init",
		Message: "",
		Value:   nil,
	})

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

	log.Printf("DisableAppInstanceKDUController - Request details for federation: %s, appInstanceId: %s, kduId: %s", federationContextId, appInstanceId, request.KduId)

	// get the appInstance from the database
	log.Printf("DisableAppInstanceKDUController - Getting app instance from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("DisableAppInstanceKDUController - Error getting app instance from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// get the application from the database
	log.Printf("DisableAppInstanceKDUController - Getting application from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	application, err := amc.applicationService.GetApplication(federationContextId, appInstance.AppId)
	if err != nil {
		log.Printf("DisableAppInstanceKDUController - Error getting application from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application: "+err.Error())
		return
	}

	// get the appPkg from the orchestrator database
	log.Printf("DisableAppInstanceKDUController - Getting app package from orchestrator for federation: %s, appInstanceId: %s, appPkgId: %s", federationContextId, appInstanceId, application.AppPkgId)
	appPkg, err := amc.orchestratorService.GetAppPkg(application.AppPkgId)
	if err != nil {
		log.Printf("DisableAppInstanceKDUController - Error getting app package from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application package: "+err.Error())
		return
	}

	err = amc.orchestratorService.DisableAppiKDU(appPkg.AppdId, request.KduId, request.NsId)
	if err != nil {
		log.Printf("DisableAppInstanceKDUController - Error disabling KDU for federation %s, appInstanceId %s, kduId %s: %v", federationContextId, appInstanceId, request.KduId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error disabling application instance KDU: "+err.Error())
		return
	}

	log.Printf("DisableAppInstanceKDUController - KDU disabled successfully for federation: %s, appInstanceId: %s, kduId: %s", federationContextId, appInstanceId, request.KduId)

	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-disable-kdu-done",
		Message: "",
		Value:   nil,
	})

	c.JSON(http.StatusOK, gin.H{"appInstance": appInstance, "kduId": request.KduId, "nsId": request.NsId})
}

// @Summary Migrate Application Instance Node
// @Description Migrates a Kubernetes Deployment Unit (KDU) within an application instance to a different node
// @Tags EWBI - ApplicationInstanceLifecycleManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param appInstanceId path string true "Application Instance ID"
// @Param request body dto.AppInstanceNodeMigrateRequest true "Node migration request"
// @Success 200 {object} map[string]string "appInstanceId: id"
// @Failure 400 {object} models.ProblemDetails "Invalid request body"
// @Failure 500 {object} models.ProblemDetails "Internal Server Error"
// @Router /ewbi/{federationContextId}/app_instances/{appInstanceId}/node/migrate [post]
func (amc *ApplicationInstanceLifecycleManagementController) AppInstanceNodeMigrateController(c *gin.Context) {
	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-migrate-appi-init",
		Message: "",
		Value:   nil,
	})

	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	// get the appInstanceId from the path
	appInstanceId := c.Param("appInstanceId")

	log.Printf("AppInstanceNodeMigrateController - Starting node migration for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)

	// get the rest of the details from the request body
	log.Printf("AppInstanceNodeMigrateController - Binding request body for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	var request dto.AppInstanceNodeMigrateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("AppInstanceNodeMigrateController - Error binding request body for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("AppInstanceNodeMigrateController - Request details for federation: %s, appInstanceId: %s, kduId: %s, nodeId: %s", federationContextId, appInstanceId, request.KduId, request.Node)

	// get the appInstance from the database
	log.Printf("AppInstanceNodeMigrateController - Getting app instance from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	appInstance, err := amc.appInstanceService.GetAppInstance(federationContextId, appInstanceId)
	if err != nil {
		log.Printf("AppInstanceNodeMigrateController - Error getting app instance from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application instance: "+err.Error())
		return
	}

	// check for the appi in the orchestrator database
	log.Printf("AppInstanceNodeMigrateController - Getting appi from orchestrator for federation: %s, appInstanceId: %s, appiId: %s", federationContextId, appInstanceId, appInstance.AppiId)
	_, err = amc.orchestratorService.GetAppi(appInstance.AppiId)
	if err != nil {
		log.Printf("AppInstanceNodeMigrateController - Error getting appi from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application package: "+err.Error())
		return
	}

	// get the application from the database
	log.Printf("AppInstanceNodeMigrateController - Getting application from database for federation: %s, appInstanceId: %s", federationContextId, appInstanceId)
	application, err := amc.applicationService.GetApplication(federationContextId, appInstance.AppId)
	if err != nil {
		log.Printf("AppInstanceNodeMigrateController - Error getting application from database for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application: "+err.Error())
		return
	}

	// get the appPkg from the orchestrator database
	log.Printf("AppInstanceNodeMigrateController - Getting app package from orchestrator for federation: %s, appInstanceId: %s, appPkgId: %s", federationContextId, appInstanceId, application.AppPkgId)
	appPkg, err := amc.orchestratorService.GetAppPkg(application.AppPkgId)
	if err != nil {
		log.Printf("AppInstanceNodeMigrateController - Error getting app package from orchestrator for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting application package: "+err.Error())
		return
	}

	// send the migrate node request to the partner
	log.Printf("AppInstanceNodeMigrateController - Sending migrate node request to partner for federation: %s, appInstanceId: %s, nsId: %s, vnfId: %s, node: %s", federationContextId, appInstanceId, request.NsId, request.VnfId, request.Node)
	err = amc.orchestratorService.MigrateAppiNode(appPkg.AppdId, request.NsId, request.VnfId, request.KduId, request.Node)
	if err != nil {
		log.Printf("AppInstanceNodeMigrateController - Error sending migrate node request to partner for federation %s, appInstanceId %s: %v", federationContextId, appInstanceId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error sending migrate node request to partner: "+err.Error())
		return
	}

	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-migrate-appi-done",
		Message: "",
		Value:   nil,
	})

	c.JSON(http.StatusOK, gin.H{"appInstanceId": appInstanceId})
}
