package models

type DynamicQoSCap struct {
	CapabilityId *CapabilityId `json:"capabilityId"`
	SupportedQoS string        `json:"supportedQoS"`
}
