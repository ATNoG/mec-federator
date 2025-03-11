package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/services"
)

// middleware to verify if the access token is valid
func AuthMiddleware(as *services.AuthServiceImpl) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			var problemDetails = &models.ProblemDetails{
				Title:  "Unauthorized",
				Detail: "Authorization Header is missing",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, problemDetails)
			return
		}

		log.Print(authorizationHeader)

		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 {
			var problemDetails = &models.ProblemDetails{
				Title:  "Unauthorized",
				Detail: "Invalid Authorization Header",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, problemDetails)
			return
		}

		tokenStr := parts[1]

		_, err := as.QueryAccessToken(tokenStr)
		if err != nil {
			var problemDetails = &models.ProblemDetails{
				Title:  "Unauthorized",
				Detail: "Provided access token is invalid or expired",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, problemDetails)
			return
		}
	}
}
