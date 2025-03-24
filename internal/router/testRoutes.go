package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/controller"
)

// TestAPIManagement - Test API for testing the service
func initTestRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	TestAPIManagement := router.Group("/test")
	TestAPIManagement.GET("/ping", controller.PingController)
	TestAPIManagement.GET("/auth", authMiddleware, controller.PingController)
}
