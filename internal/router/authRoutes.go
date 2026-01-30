package router

import (
	"github.com/atnog/mec-federator/internal/config"
	"github.com/atnog/mec-federator/internal/controller"
	"github.com/gin-gonic/gin"
)

// AuthAPIManagement - Authentication and authorization of the partner OP
func initAuthRoutes(router *gin.Engine, svcs *Services) {
	AuthAPIManagement := router.Group("/federation/v1/auth")

	authController := controller.NewAuthController(
		svcs.AuthService,
		config.GetMongoClient(),
	)

	AuthAPIManagement.POST(
		"/token",
		authController.IssueAccessTokenController)
}
