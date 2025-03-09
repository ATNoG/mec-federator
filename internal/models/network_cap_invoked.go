package models

import (
	"time"
)

type NetworkCapInvoked struct {
	NetworkEventId string `json:"networkEventId"`

	CapabilityId *CapabilityId `json:"capabilityId"`

	InvocationTime *time.Time `json:"invocationTime,omitempty"`

	NwCapabilitySLI string `json:"nwCapabilitySLI"`

	ZoneId string `json:"zoneId"`
}
