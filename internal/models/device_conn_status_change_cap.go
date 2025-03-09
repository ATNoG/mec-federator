package models

type DeviceConnStatusChangeCap struct {
	CapabilityId      *CapabilityId `json:"capabilityId"`
	MaxiDetectionTime string        `json:"maxiDetectionTime"`
}
