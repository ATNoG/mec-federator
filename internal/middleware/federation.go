package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/services"
)

// middleware to check if the federation exists
func FederationMiddleware(federationService *services.FederationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		federationId := c.Param("federationId")
		_, err := federationService.GetFederation(federationId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Federation not found"})
		}
	}
}
