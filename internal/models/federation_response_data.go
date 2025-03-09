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

	PartnerOPFixedNetworkCodes *[]string `json:"partnerOPFixedNetworkCodes,omitempty"`
	// List of zones, which the operator platform wishes to make available to developers/ISVs of requesting operator platform.
	OfferedAvailabilityZones []ZoneDetails `json:"offeredAvailabilityZones,omitempty"`

	PlatformCaps *[]string `json:"platformCaps"`
	// Date and Time zone info of the existing federation expiry
	FederationExpiryDate time.Time `json:"federationExpiryDate,omitempty"`
	// Date and Time zone info of the existing federation renewal. Shall be less than federationExpiryDate
	FederationRenewalDate time.Time `json:"federationRenewalDate,omitempty"`
}
