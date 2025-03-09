package models

type UserPlaneMgmtEvtCap struct {
	CapabilityId        *CapabilityId `json:"capabilityId"`
	MaxUserPlaneLatency string        `json:"maxUserPlaneLatency"`
}
