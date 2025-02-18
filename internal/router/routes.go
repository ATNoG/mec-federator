package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/handler"
)

func initRoutes(router *gin.Engine) {
	v1 := router.Group("/operatorplatform/federation/v1")
	{
		v1.POST("/test", handler.TestHandler)

	}
}
