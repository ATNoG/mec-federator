package utils

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
)

func HandleProblem(c *gin.Context, code int, detail string) {
	log.Printf("Handling problem: %d, %s", code, detail)
	problemDetails := NewProblemDetails(code, detail)
	c.AbortWithStatusJSON(code, problemDetails)
}

func NewProblemDetails(code int, detail string) *models.ProblemDetails {
	var problemDetails models.ProblemDetails

	switch code {
	case 400:
		problemDetails.Title = "Bad Request."
	case 401:
		problemDetails.Title = "Unauthorized."
	case 404:
		problemDetails.Title = "Content Not Found."
	case 409:
		problemDetails.Title = "Conflict."
	case 422:
		problemDetails.Title = "Unprocessable Entity."
	case 500:
		problemDetails.Title = "Internal Server Error."
	case 503:
		problemDetails.Title = "Service Unavailable."

	default:
		problemDetails.Title = "Unknown Error."
	}

	if detail != "" {
		problemDetails.Detail = detail
	}

	return &problemDetails
}
