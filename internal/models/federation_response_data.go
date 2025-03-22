package models

import (
	"time"
)

type FederationResponseData struct {
	PartnerOPFederationId string `json:"partnerOPFederationId,omitempty" bson:"partnerOPFederationId,omitempty"`

	PartnerOPCountryCode string `json:"partnerOPCountryCode" bson:"partnerOPCountryCode" binding:"required"`

	FederationContextId string `json:"federationContextId" bson:"federationContextId" binding:"required"`

	EdgeDiscoveryServiceEndPoint *ServiceEndpoint `json:"edgeDiscoveryServiceEndPoint" bson:"edgeDiscoveryServiceEndPoint" binding:"required"`

	LcmServiceEndPoint *ServiceEndpoint `json:"lcmServiceEndPoint" bson:"lcmServiceEndPoint" binding:"required"`

	PartnerOPMobileNetworkCodes *MobileNetworkIds `json:"partnerOPMobileNetworkCodes,omitempty" bson:"partnerOPMobileNetworkCodes,omitempty"`

	PartnerOPFixedNetworkCodes *[]string `json:"partnerOPFixedNetworkCodes,omitempty" bson:"partnerOPFixedNetworkCodes,omitempty"`
	// List of zones, which the operator platform wishes to make available to developers/ISVs of requesting operator platform.
	OfferedAvailabilityZones []ZoneDetails `json:"offeredAvailabilityZones,omitempty" bson:"offeredAvailabilityZones,omitempty"`

	PlatformCaps *[]string `json:"platformCaps" bson:"platformCaps" binding:"required"`
	// Date and Time zone info of the existing federation expiry
	FederationExpiryDate time.Time `json:"federationExpiryDate,omitempty" bson:"federationExpiryDate,omitempty"`
	// Date and Time zone info of the existing federation renewal. Shall be less than federationExpiryDate
	FederationRenewalDate time.Time `json:"federationRenewalDate,omitempty" bson:"federationRenewalDate,omitempty"`
}
