package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BeginAuthController(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "Begin Auth"})
}
