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

func initZoneInfoSyncRoutes(router *gin.Engine, svcs *Services, mdws *Middlewares) {
	// ZoneInfoSync - Sync zone information
	ZoneInfoSync := router.Group("/federation/v1/ewbi", *mdws.AuthMiddleware)

	zoneInfoSyncController := ewbi.NewZonesInfoSyncController(
		svcs.ZoneService,
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
