package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initRoutes(router *gin.Engine) {
	v1 := router.Group("/operatorplatform/federation/v1")
	{
		v1.POST("/partner", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "POST /operatorplatform/federation/v1/partner",
			})
		})
	}
}
