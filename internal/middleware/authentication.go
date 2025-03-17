package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

// middleware to verify if the access token is valid
func AuthMiddleware(as *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			problemDetails := utils.NewProblemDetails(http.StatusUnauthorized)
			problemDetails.Detail = "Authorization Header is missing"
			c.AbortWithStatusJSON(http.StatusUnauthorized, problemDetails)
			return
		}

		log.Print("AuthMiddleware - Found Authorization Header")

		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 {
			problemDetails := utils.NewProblemDetails(http.StatusUnauthorized)
			problemDetails.Detail = "Invalid Authorization Header"
			c.AbortWithStatusJSON(http.StatusUnauthorized, problemDetails)
			return
		}

		tokenStr := parts[1]
		_, err := as.QueryAccessToken(tokenStr)
		if err != nil {
			problemDetails := utils.NewProblemDetails(http.StatusUnauthorized)
			problemDetails.Detail = "Provided access token is invalid or expired"
			c.AbortWithStatusJSON(http.StatusUnauthorized, problemDetails)
			return
		}
	}
}
