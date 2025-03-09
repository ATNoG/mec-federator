package models

type CapabilityId string

const (
	CONN_STATE_CHANGE     CapabilityId = "NW_CAP_CONN_STATE_CHANGE"
	LOCATION_RETRIEVAL    CapabilityId = "NW_CAP_LOCATION_RETRIEVAL"
	USERPLANE_MGMT_EVENTS CapabilityId = "NW_CAP_USERPLANE_MGMT_EVENTS"
	DYNAMIC_QOS           CapabilityId = "NW_CAP_DYNAMIC_QOS"
)
