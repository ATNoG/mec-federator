package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

// middleware to check if the federation exists
func FederationExistsMiddleware(federationService *services.FederationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		federationId := c.Param("federationId")
		_, err := federationService.GetFederation(federationId)
		if err != nil {
			log.Print("FederationMiddleware - Federation not found")
			utils.HandleProblem(c, http.StatusNotFound, "Federation not found")
			return
		}
	}
}
