package models

import (
	"time"
)

type FederationHealthInfo struct {
	FederationStatus *State `json:"federationStatus"`

	FederationStartTime *time.Time `json:"federationStartTime"`

	NumOfAcceptedZones string `json:"numOfAcceptedZones"`

	NumOfActiveAlarms string `json:"numOfActiveAlarms,omitempty"`

	NumOfApplications string `json:"numOfApplications,omitempty"`
}
