package models

type DeviceConnStatusChangeCap struct {
	CapabilityId *CapabilityId `json:"capabilityId"`
	// The maximum detection time in seconds that the OP can determine the UE change of connectivity with the mobile network.
	MaxiDetectionTime string `json:"maxiDetectionTime"`
}
