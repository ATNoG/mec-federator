package models

type InvalidParam struct {
	Param string `json:"param"`

	Reason string `json:"reason,omitempty"`
}
