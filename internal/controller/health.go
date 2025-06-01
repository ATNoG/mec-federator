package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (hc *HealthController) HealthCheckController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Healthy."})
}
