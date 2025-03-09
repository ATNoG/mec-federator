package models

// Application metadata details
type AppMetaData struct {
	// Name of the application. Application provider define a human readable name for the application
	AppName string `json:"appName"`
	// Version info of the application
	Version string `json:"version"`
	// Brief application description provided by application provider
	AppDescription string `json:"appDescription,omitempty"`

	MobilitySupport bool `json:"mobilitySupport,omitempty"`
	// An application Access key, to be used with UNI interface to authorize UCs Access to a given application
	AccessToken string `json:"accessToken"`
	// Possible categorization of the application
	Category string `json:"category,omitempty"`
}
