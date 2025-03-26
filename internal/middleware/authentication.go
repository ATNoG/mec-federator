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
			log.Print("AuthMiddleware - Authorization Header is missing")
			utils.HandleProblem(c, http.StatusUnauthorized, "Authorization Header is missing")
			return
		}

		log.Print("AuthMiddleware - Found Authorization Header")

		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 {
			log.Print("AuthMiddleware - Invalid Authorization Header")
			utils.HandleProblem(c, http.StatusUnauthorized, "Invalid Authorization Header")
			return
		}

		tokenStr := parts[1]
		_, err := as.QueryAccessToken(tokenStr)
		if err != nil {
			log.Print("AuthMiddleware - Provided access token is invalid or expired")
			utils.HandleProblem(c, http.StatusUnauthorized, "Provided access token is invalid or expired")
			return
		}

		log.Print("AuthMiddleware - Access token is valid, setting it in the context")
		c.Set("accessToken", tokenStr)
	}
}
