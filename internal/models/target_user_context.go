package models

type TargetUserContext struct {
	ConnectID string `json:"connectID"`

	ExpiryDuration *ExpiryInterval `json:"expiryDuration"`
}
