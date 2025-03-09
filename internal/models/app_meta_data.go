package models

type AppMetaData struct {
	AppName        string `json:"appName"`
	Version        string `json:"version"`
	AppDescription string `json:"appDescription,omitempty"`

	MobilitySupport bool   `json:"mobilitySupport,omitempty"`
	AccessToken     string `json:"accessToken"`
	Category        string `json:"category,omitempty"`
}
