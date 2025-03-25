package models

import (
	"time"
)

type FederationHealthInfo struct {
	FederationStatus *State `json:"federationStatus" bson:"federationStatus"`

	FederationStartTime time.Time `json:"federationStartTime" bson:"federationStartTime"`

	NumOfAcceptedZones string `json:"numOfAcceptedZones" bson:"numOfAcceptedZones"`

	NumOfActiveAlarms string `json:"numOfActiveAlarms,omitempty" bson:"numOfActiveAlarms"`

	NumOfApplications string `json:"numOfApplications,omitempty" bson:"numOfApplications"`
}
