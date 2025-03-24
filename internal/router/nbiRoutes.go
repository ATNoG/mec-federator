package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller/nbi"
)

func initNbiFederationManagementRoutes(router *gin.Engine, svcs *Services) {
	federationManagementController := nbi.NewFederationManagementController(svcs.FederationService, svcs.FederationHttpClientManager)
	// FederationManagement - Create and manage directed federation relationship with a partner OP
	FederationManagement := router.Group("/nbi/federation/v1")

	FederationManagement.POST("/partner", federationManagementController.InitiateFederationController)
	// FederationManagement.GET("/fed-context-id", controller.GetFederationContextId)
	// FederationManagement.GET("/:federationContextId/health", controller.GetFederationHealth)
	// FederationManagement.POST("/:federationContextId/renew", controller.PostFederationRenew)
	// FederationManagement.GET("/:federationContextId/platform-caps", controller.GetFederationPlatformCaps)
}
