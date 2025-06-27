package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller/ewbi"
)

func initEwbiFederationManagementRoutes(router *gin.Engine, svcs *Services, mdws *Middlewares) {
	// FederationManagement - Create and manage directed federation relationship with a partner OP
	FederationManagement := router.Group("/federation/v1/ewbi", *mdws.AuthMiddleware)

	federationManagementController := ewbi.NewFederationManagementController(
		svcs.FederationService,
		svcs.ZoneService,
		svcs.ArtefactService,
		svcs.AppInstanceService,
	)

	FederationManagement.POST(
		"/partner",
		federationManagementController.CreateFederationController)
	FederationManagement.GET(
		"/:federationContextId/partner",
		*mdws.FederationExistsMiddleware,
		federationManagementController.GetFederationMetaInfoController)
	FederationManagement.PATCH(
		"/:federationContextId/partner",
		*mdws.FederationExistsMiddleware,
		federationManagementController.UpdateFederationController)
	FederationManagement.DELETE(
		"/:federationContextId/partner",
		*mdws.FederationExistsMiddleware,
		federationManagementController.RemoveFederationController)
	FederationManagement.GET(
		"/fed-context-id",
		federationManagementController.GetFederationContextIdentifierController)
	FederationManagement.GET(
		"/:federationContextId/health",
		*mdws.FederationExistsMiddleware,
		federationManagementController.GetFederationHealthController)
	FederationManagement.POST(
		"/:federationContextId/renew",
		*mdws.FederationExistsMiddleware,
		federationManagementController.RenewFederationController)
}

func initEwbiZoneInfoSyncRoutes(router *gin.Engine, svcs *Services, mdws *Middlewares) {
	// ZoneInfoSync - Sync zone information
	ZoneInfoSync := router.Group("/federation/v1/ewbi", *mdws.AuthMiddleware)

	zoneInfoSyncController := ewbi.NewZonesInfoSyncController(
		svcs.ZoneService,
		svcs.OrchestratorService,
	)

	ZoneInfoSync.POST(
		"/:federationContextId/zones",
		*mdws.FederationExistsMiddleware,
		zoneInfoSyncController.SubscribeZoneController)
	ZoneInfoSync.DELETE(
		"/:federationContextId/zones/:zoneId",
		*mdws.FederationExistsMiddleware,
		zoneInfoSyncController.UnsubscribeZoneController)
	ZoneInfoSync.GET(
		"/:federationContextId/zones/:zoneId",
		*mdws.FederationExistsMiddleware,
		zoneInfoSyncController.GetZoneController)

	ZoneInfoSync.GET(
		"/:federationContextId/zones",
		*mdws.FederationExistsMiddleware,
		zoneInfoSyncController.GetAllLocalZonesController)
}

func initEwbiArtefactManagementRoutes(router *gin.Engine, svcs *Services, mdws *Middlewares) {
	// ArtefactManagement - Create and manage artefacts
	ArtefactManagement := router.Group("/federation/v1/ewbi", *mdws.AuthMiddleware)

	artefactManagementController := ewbi.NewArtefactManagementController(
		svcs.OrchestratorService,
		svcs.ArtefactService,
	)

	ArtefactManagement.POST(
		"/:federationContextId/artefact",
		*mdws.FederationExistsMiddleware,
		artefactManagementController.OnboardArtefactController)
	ArtefactManagement.GET(
		"/:federationContextId/artefact/:artefactId",
		*mdws.FederationExistsMiddleware,
		artefactManagementController.GetArtefactController)
	ArtefactManagement.DELETE(
		"/:federationContextId/artefact/:artefactId",
		*mdws.FederationExistsMiddleware,
		artefactManagementController.DeleteArtefactController)
	ArtefactManagement.POST(
		"/:federationContextId/files",
		*mdws.FederationExistsMiddleware,
		artefactManagementController.UploadFileController)
	ArtefactManagement.GET(
		"/:federationContextId/files/:fileId",
		*mdws.FederationExistsMiddleware,
		artefactManagementController.GetFileController)
	ArtefactManagement.DELETE(
		"/:federationContextId/files/:fileId",
		*mdws.FederationExistsMiddleware,
		artefactManagementController.DeleteFileController)
}

func initEwbiApplicationInstanceLifecycleManagementRoutes(router *gin.Engine, svcs *Services, mdws *Middlewares) {
	// ApplicationInstanceLifecycleManagement - Create and manage application instance lifecycle
	ApplicationInstanceLifecycleManagement := router.Group("/federation/v1/ewbi", *mdws.AuthMiddleware)

	applicationInstanceLifecycleManagementController := ewbi.NewApplicationInstanceLifecycleManagementController(
		svcs.OrchestratorService,
		svcs.ArtefactService,
		svcs.AppInstanceService,
		svcs.ZoneService,
	)

	ApplicationInstanceLifecycleManagement.POST(
		"/:federationContextId/application/lcm",
		*mdws.FederationExistsMiddleware,
		applicationInstanceLifecycleManagementController.CreateAppInstanceController)
	ApplicationInstanceLifecycleManagement.DELETE(
		"/:federationContextId/application/lcm/:appInstanceId",
		*mdws.FederationExistsMiddleware,
		applicationInstanceLifecycleManagementController.DeleteAppInstanceController)
	ApplicationInstanceLifecycleManagement.GET(
		"/:federationContextId/application/lcm/:appInstanceId",
		*mdws.FederationExistsMiddleware,
		applicationInstanceLifecycleManagementController.GetAppInstanceDetailsController)

	ApplicationInstanceLifecycleManagement.PATCH(
		"/:federationContextId/application/lcm/:appInstanceId/kdu/enable",
		*mdws.FederationExistsMiddleware,
		applicationInstanceLifecycleManagementController.EnableAppInstanceKDUController)
	ApplicationInstanceLifecycleManagement.PATCH(
		"/:federationContextId/application/lcm/:appInstanceId/kdu/disable",
		*mdws.FederationExistsMiddleware,
		applicationInstanceLifecycleManagementController.DisableAppInstanceKDUController)
}
