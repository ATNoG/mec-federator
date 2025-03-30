package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller/ewbi"
)

func initEwbiFederationManagementRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	// FederationManagement - Create and manage directed federation relationship with a partner OP
	FederationManagement := router.Group("/federation/v1/ewbi", authMiddleware)

	federationManagementController := ewbi.NewFederationManagementController(
		svcs.FederationService,
	)

	FederationManagement.POST(
		"/partner",
		federationManagementController.CreateFederationController)
	FederationManagement.GET(
		"/:federationContextId/partner",
		federationManagementController.GetFederationMetaInfoController)
	FederationManagement.PATCH(
		"/:federationContextId/partner",
		federationManagementController.UpdateFederationController)
	FederationManagement.DELETE(
		"/:federationContextId/partner",
		federationManagementController.RemoveFederationController)
	FederationManagement.GET(
		"/fed-context-id",
		federationManagementController.GetFederationContextIdentifierController)
	FederationManagement.GET(
		"/:federationContextId/health",
		federationManagementController.GetFederationHealthController)
	FederationManagement.POST(
		"/:federationContextId/renew",
		federationManagementController.RenewFederationController)
}

func initEwbiArtefactManagementRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	// ArtefactManagement - Create and manage artefacts
	ArtefactManagement := router.Group("/artefacts/v1/ewbi", authMiddleware)

	artefactManagementController := ewbi.NewArtefactManagementController(
		svcs.OrchestratorService,
	)

	ArtefactManagement.POST(
		"/:federationContextId/artefact",
		artefactManagementController.OnboardArtefactController)
	ArtefactManagement.GET(
		"/:federationContextId/artefact/:artefactId",
		artefactManagementController.GetArtefactController)
	ArtefactManagement.DELETE(
		"/:federationContextId/artefact/:artefactId",
		artefactManagementController.DeleteArtefactController)
	ArtefactManagement.POST(
		"/:federationContextId/files",
		artefactManagementController.UploadFileController)
	ArtefactManagement.GET(
		"/:federationContextId/files/:fileId",
		artefactManagementController.GetFileController)
	ArtefactManagement.DELETE(
		"/:federationContextId/files/:fileId",
		artefactManagementController.DeleteFileController)
}
