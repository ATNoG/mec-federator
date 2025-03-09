package models

type ZoneRegisteredData struct {
	ZoneId                     string                `json:"zoneId"`
	ReservedComputeResources   []ComputeResourceInfo `json:"reservedComputeResources"`
	ComputeResourceQuotaLimits []ComputeResourceInfo `json:"computeResourceQuotaLimits"`

	FlavoursSupported []Flavour `json:"flavoursSupported"`

	NetworkResources         *interface{} `json:"networkResources,omitempty"`
	ZoneServiceLevelObjsInfo *interface{} `json:"zoneServiceLevelObjsInfo,omitempty"`
}
