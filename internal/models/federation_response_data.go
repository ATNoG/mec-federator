package models

import (
	"time"
)

type FederationResponseData struct {
	PartnerOPFederationId string `json:"partnerOPFederationId,omitempty"`

	PartnerOPCountryCode string `json:"partnerOPCountryCode,omitempty"`

	FederationContextId string `json:"federationContextId"`

	EdgeDiscoveryServiceEndPoint *ServiceEndpoint `json:"edgeDiscoveryServiceEndPoint,omitempty"`

	LcmServiceEndPoint *ServiceEndpoint `json:"lcmServiceEndPoint,omitempty"`

	PartnerOPMobileNetworkCodes *MobileNetworkIds `json:"partnerOPMobileNetworkCodes,omitempty"`

	PartnerOPFixedNetworkCodes *[]string     `json:"partnerOPFixedNetworkCodes,omitempty"`
	OfferedAvailabilityZones   []ZoneDetails `json:"offeredAvailabilityZones,omitempty"`

	PlatformCaps          *[]string `json:"platformCaps"`
	FederationExpiryDate  time.Time `json:"federationExpiryDate,omitempty"`
	FederationRenewalDate time.Time `json:"federationRenewalDate,omitempty"`
}
