package models

type ProblemDetails struct {
	// Summary of the problem
	Title string `json:"title,omitempty"`
	// Specific detail of the issue
	Detail string `json:"detail,omitempty"`
	// Fixed string indicating cause of the issue
	Cause string `json:"cause,omitempty"`

	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}
