package models

type UserPlaneMgmtEvtCap struct {
	CapabilityId *CapabilityId `json:"capabilityId"`
	// Indicates the maximum user plane latency in units of milliseconds to decide whether edge relocation is needed to ascertain latency remain in this range.
	MaxUserPlaneLatency string `json:"maxUserPlaneLatency"`
}
