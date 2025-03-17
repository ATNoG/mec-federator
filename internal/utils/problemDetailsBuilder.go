package utils

import "github.com/mankings/mec-federator/internal/models"

func NewProblemDetails(code int) *models.ProblemDetails {
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

	return &problemDetails
}
