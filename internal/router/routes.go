package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/controller"
	"github.com/mankings/mec-federator/internal/controller/ewbi"
)

func initRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	initTestRoutes(router, svcs, authMiddleware)
	initAuthRoutes(router, svcs, authMiddleware)

	// EWBI API Management
	initFederationAPIManagementRoutes(router, svcs, authMiddleware)
	initFederationManagementRoutes(router, svcs, authMiddleware)

	// SBI API Management
	// orchestrator registration routes
}

func initTestRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	// TestAPIManagement - Test API for testing the service
	TestAPIManagement := router.Group("/federation/v1/test")
	TestAPIManagement.GET("/ping", controller.PingController)
	TestAPIManagement.GET("/health", controller.HealthController)
	TestAPIManagement.GET("/auth", authMiddleware, controller.HealthController)
}

func initAuthRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	authController := controller.NewAuthController(svcs.AuthService, config.GetMongoClient())
	// AuthAPIManagement - Authentication and authorization of the partner OP
	AuthAPIManagement := router.Group("/auth")
	AuthAPIManagement.POST("/token", authController.IssueAccessTokenController)
}

func initFederationAPIManagementRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	// FederationAPIManagement - Retrieves federation resources and methods a partner OP support on E/WBI
	// FederationAPIManagement := router.Group("/federation/v1")

	// FederationAPIManagement.GET("/federation-resources", controller.GetFederationResourcesController)
}

func initFederationManagementRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	federationManagementController := ewbi.NewFederationManagementController(svcs.FederationService)
	// FederationManagement - Create and manage directed federation relationship with a partner OP
	FederationManagement := router.Group("/federation/v1", authMiddleware)

	FederationManagement.POST("/partner", federationManagementController.CreateFederationController)
	FederationManagement.DELETE("/:federationContextId/partner", federationManagementController.RemoveFederationRelationshipController)
	// FederationManagement.GET("/:federationContextId/partner", controller.GetFederationPartnerController)
	// FederationManagement.PATCH("/:federationContextId/partner", controller.PatchFederationPartnerController)
	// FederationManagement.GET("/fed-context-id", controller.GetFederationContextId)
	// FederationManagement.GET("/:federationContextId/health", controller.GetFederationHealth)
	// FederationManagement.POST("/:federationContextId/renew", controller.PostFederationRenew)
	// FederationManagement.GET("/:federationContextId/platform-caps", controller.GetFederationPlatformCaps)
}
