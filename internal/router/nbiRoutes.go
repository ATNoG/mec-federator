package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller/nbi"
)

func initNbiFederationManagementRoutes(router *gin.Engine, svcs *Services) {
	// FederationManagement - Create and manage directed federation relationship with a partner OP
	FederationManagement := router.Group("/federation/v1/nbi")

	federationManagementController := nbi.NewFederationManagementController(
		svcs.FederationService,
		svcs.HttpClientService,
	)

	FederationManagement.POST(
		"/partner",
		federationManagementController.InitiateFederationController)
	FederationManagement.GET(
		"/:federationContextId/partner",
		federationManagementController.GetFederationMetaInfoController)
	FederationManagement.DELETE(
		"/:federationContextId/partner",
		federationManagementController.RemoveFederationController)
	FederationManagement.GET(
		"/:federationContextId/health",
		federationManagementController.GetFederationHealthController)
	FederationManagement.GET(
		"/:federationContextId/renew",
		federationManagementController.RequestFederationRenewalController)
}
