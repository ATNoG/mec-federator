package models

type ProblemDetails struct {
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Cause  string `json:"cause,omitempty"`

	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}
