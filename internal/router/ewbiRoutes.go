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
	FederationManagement.DELETE(
		"/:federationContextId/partner",
		federationManagementController.RemoveFederationController)
	FederationManagement.GET(
		"/:federationContextId/partner",
		federationManagementController.GetFederationMetaInfoController)
	FederationManagement.PATCH(
		"/:federationContextId/partner",
		federationManagementController.UpdateFederationController)
	FederationManagement.GET(
		"/:federationContextId/health",
		federationManagementController.GetFederationHealthController)
	FederationManagement.POST(
		"/:federationContextId/renew",
		federationManagementController.RenewFederationController)
}
