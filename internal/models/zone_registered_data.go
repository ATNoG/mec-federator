package models

type ZoneRegisteredData struct {
	ZoneId string `json:"zoneId" bson:"zoneId"`

	FederationContextId string `json:"federationContextId" bson:"federationContextId"`

	// Resources exclusively reserved for the originator OP.
	ReservedComputeResources []ComputeResourceInfo `json:"reservedComputeResources" bson:"reservedComputeResources"`
	// Max quota on resources partner OP allows over reserved resources.
	ComputeResourceQuotaLimits []ComputeResourceInfo `json:"computeResourceQuotaLimits" bson:"computeResourceQuotaLimits"`

	FlavoursSupported []Flavour `json:"flavoursSupported" bson:"flavoursSupported"`

	NetworkResources *interface{} `json:"networkResources,omitempty" bson:"networkResources,omitempty"`
	// It is a measure of the actual amount of data that is being sent over a network per unit of time and indicates máximum supported value for a zone
	ZoneServiceLevelObjsInfo *interface{} `json:"zoneServiceLevelObjsInfo,omitempty" bson:"zoneServiceLevelObjsInfo,omitempty"`
}
