package models

type FederationMetaInfo struct {
	EdgeDiscoveryServiceEndPoint *ServiceEndpoint  `json:"edgeDiscoveryServiceEndPoint"`
	OfferedAvailabilityZones     []ZoneDetails     `json:"offeredAvailabilityZones"`
	AllowedMobileNetworkIds      *MobileNetworkIds `json:"allowedMobileNetworkIds"`
	AllowedFixedNetworkIds       *[]string         `json:"allowedFixedNetworkIds"`
	LcmServiceEndPoint           *ServiceEndpoint  `json:"lcmServiceEndPoint"`
	PlatformCaps                 *[]string         `json:"platformCaps" binding:"required"`
}
