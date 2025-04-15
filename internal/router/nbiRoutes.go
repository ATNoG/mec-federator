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

func initMecSystemManagementRoutes(router *gin.Engine, svcs *Services) {
	// OrchestratorManagement - Register and manage orchestrator information
	OrchestratorManagement := router.Group("/federation/v1/nbi/mec-system")

	orchestratorManagementController := nbi.NewMecSystemManagementController(
		svcs.OrchestratorService,
		svcs.MecSystemService,
	)

	OrchestratorManagement.POST(
		"",
		orchestratorManagementController.RegisterOrchestratorController)
	OrchestratorManagement.GET(
		"",
		orchestratorManagementController.GetOrchestratorInfoController)
}
