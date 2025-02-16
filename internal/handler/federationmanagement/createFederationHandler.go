package federationmanagement

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateFederationHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "POST /operatorplatform/federation/v1/partner",
	})
}
