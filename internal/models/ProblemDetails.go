package models

// ProblemDetails is a struct that represents the problem details object
type ProblemDetails struct {
	Title         string          `json:"title"`
	Details       string          `json:"details"`
	Cause         string          `json:"cause"`
	InvalidParams []InvalidParams `json:"invalidParams"`
}
