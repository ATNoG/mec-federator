package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
)

func GetAuthenticationToken(c *gin.Context) string {
	authorizationHeader := c.GetHeader("Authorization")
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 {
		var problemDetails = &models.ProblemDetails{
			Title:  "Unauthorized",
			Detail: "Invalid Authorization Header",
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, problemDetails)
	}
	return parts[1]
}
