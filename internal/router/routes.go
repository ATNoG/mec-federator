package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller"
)

func initRoutes(router *gin.Engine) {
	initTestRoutes(router)
	initAuthRoutes(router)
	initFederationAPIManagementRoutes(router)
	initFederationManagementRoutes(router)
}

func initTestRoutes(router *gin.Engine) {
	// TestAPIManagement - Test API for testing the service
	TestAPIManagement := router.Group("/test")
	TestAPIManagement.GET("/ping", controller.PingController)
	TestAPIManagement.GET("/health", controller.HealthController)
}

func initAuthRoutes(router *gin.Engine) {
	// AuthAPIManagement - Authentication and authorization of the partner OP
	AuthAPIManagement := router.Group("/auth")
	AuthAPIManagement.GET("/", controller.BeginAuthController)
	// AuthAPIManagement.GET("/token", controller.CompleteAuthController)
	// AuthAPIManagement.GET("/logout", controller.LogoutAuthController)
}

func initFederationAPIManagementRoutes(router *gin.Engine) {
	// FederationAPIManagement - Retrieves federation resources and methods a partner OP support on E/WBI
	// FederationAPIManagement := router.Group("/federation/v1")
	// FederationAPIManagement.GET("/federation-resources", controller.GetFederationResourcesController)
}

func initFederationManagementRoutes(router *gin.Engine) {
	// FederationManagement - Create and manage directed federation relationship with a partner OP
	// FederationManagement := router.Group("/federation/v1")
	// FederationManagement.POST("/partner", controller.PostFederationPartnerController)
	// FederationManagement.GET("/:federationContextId/partner", controller.GetFederationPartnerController)
	// FederationManagement.PATCH("/:federationContextId/partner", controller.PatchFederationPartnerController)
	// FederationManagement.DELETE("/:federationContextId/partner", controller.DeleteFederationPartnerController)
	// FederationManagement.GET("/fed-context-id", controller.GetFederationContextId)
	// FederationManagement.GET("/:federationContextId/health", controller.GetFederationHealth)
	// FederationManagement.POST("/:federationContextId/renew", controller.PostFederationRenew)
	// FederationManagement.GET("/:federationContextId/platform-caps", controller.GetFederationPlatformCaps)
}
