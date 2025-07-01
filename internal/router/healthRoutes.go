package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller"
)

// HealthAPIManagement - Health check of the partner OP
func initHealthRoutes(router *gin.Engine) {
	HealthAPIManagement := router.Group("")

	healthController := controller.NewHealthController()

	HealthAPIManagement.GET(
		"/healthcheck",
		healthController.HealthCheckController)
}
