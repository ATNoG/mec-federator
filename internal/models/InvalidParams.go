package models

// InvalidParams is a struct that represents the invalid parameters object
type InvalidParams struct {
	Param  string `json:"param" binding:"required"`
	Reason string `json:"reason"`
}
