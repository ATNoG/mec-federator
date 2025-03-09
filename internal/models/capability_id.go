package models

// CapabilityId : The enumerated list of network capabilities that an OP can use for various services via SBI-NR.
type CapabilityId string

// List of CapabilityID
const (
	CONN_STATE_CHANGE     CapabilityId = "NW_CAP_CONN_STATE_CHANGE"
	LOCATION_RETRIEVAL    CapabilityId = "NW_CAP_LOCATION_RETRIEVAL"
	USERPLANE_MGMT_EVENTS CapabilityId = "NW_CAP_USERPLANE_MGMT_EVENTS"
	DYNAMIC_QOS           CapabilityId = "NW_CAP_DYNAMIC_QOS"
)
