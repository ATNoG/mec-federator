package router

import (
	"github.com/atnog/mec-federator/internal/controller"
	"github.com/gin-gonic/gin"
)

// HealthAPIManagement - Health check of the partner OP
func initHealthRoutes(router *gin.Engine) {
	HealthAPIManagement := router.Group("")

	healthController := controller.NewHealthController()

	HealthAPIManagement.GET(
		"/healthcheck",
		healthController.HealthCheckController)
}
