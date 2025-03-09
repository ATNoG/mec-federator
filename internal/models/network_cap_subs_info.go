package models

type NetworkCapSubsInfo struct {
	AppId string `json:"appId"`

	AppProviderId string `json:"appProviderId"`

	CapabilityId *CapabilityId `json:"capabilityId"`
}
