package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/controller"
)

// AuthAPIManagement - Authentication and authorization of the partner OP
func initAuthRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	authController := controller.NewAuthController(svcs.AuthService, config.GetMongoClient())
	AuthAPIManagement := router.Group("/auth")
	AuthAPIManagement.POST("/token", authController.IssueAccessTokenController)
}
