package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller/ewbi"
)

func initFederationManagementRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	federationManagementController := ewbi.NewFederationManagementController(svcs.FederationService)
	// FederationManagement - Create and manage directed federation relationship with a partner OP
	FederationManagement := router.Group("/federation/v1", authMiddleware)

	FederationManagement.POST("/partner", federationManagementController.CreateFederationController)
	FederationManagement.DELETE("/:federationContextId/partner", federationManagementController.RemoveFederationRelationshipController)
	FederationManagement.GET("/:federationContextId/partner", federationManagementController.GetFederationMetaInfoController)
	FederationManagement.PATCH("/:federationContextId/partner", federationManagementController.UpdateFederationController)
	// FederationManagement.GET("/fed-context-id", controller.GetFederationContextId)
	// FederationManagement.GET("/:federationContextId/health", controller.GetFederationHealth)
	// FederationManagement.POST("/:federationContextId/renew", controller.PostFederationRenew)
	// FederationManagement.GET("/:federationContextId/platform-caps", controller.GetFederationPlatformCaps)
}

