package models

type DynamicQoSCap struct {
	CapabilityId *CapabilityId `json:"capabilityId"`
	// Set of one or more 5G QoS Identifier (5QI or 4G QCI) created via concatanation of Resource Type and 5QI values i.e., GBR1, GBR2, GBR65, NONGBR79 etc.
	SupportedQoS string `json:"supportedQoS"`
}
