package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller"
)

func initRoutes(router *gin.Engine) {
	FederationAPIManagement := router.Group("/operatorplatform/federation/v1")
	{
		FederationAPIManagement.GET("/federation-resources", controller.GetFederationResourcesController)
	}

	FederationManagement := router.Group("/operatorplatform/federation/v1")
	{
		FederationManagement.POST("/partner", controller.PostFederationPartnerController)
		// FederationManagement.GET("/:federationContextId/partner", controller.GetFederationPartnerController)
		// FederationManagement.PATCH("/:federationContextId/partner", controller.PatchFederationPartnerController)
		// FederationManagement.DELETE("/:federationContextId/partner", controller.DeleteFederationPartnerController)
		// FederationManagement.GET("/fed-context-id", controller.GetFederationContextId)
		// FederationManagement.GET("/:federationContextId/health", controller.GetFederationHealth)
		// FederationManagement.POST("/:federationContextId/renew", controller.PostFederationRenew)
		// FederationManagement.GET("/:federationContextId/platform-caps", controller.GetFederationPlatformCaps)
	}
}
