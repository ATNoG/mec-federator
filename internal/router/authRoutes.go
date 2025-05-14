package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/controller"
)

// AuthAPIManagement - Authentication and authorization of the partner OP
func initAuthRoutes(router *gin.Engine, svcs *Services) {
	AuthAPIManagement := router.Group("/auth")

	authController := controller.NewAuthController(
		svcs.AuthService,
		config.GetMongoClient(),
	)

	AuthAPIManagement.POST(
		"/token",
		authController.IssueAccessTokenController)
}
