package models

import "time"

type FederationResponseData struct {
	PartnerOPFederationId        string             `json:"partnerOPFederationId"`
	PartnerOPCountryCode         string             `json:"partnerOPCountryCode" binding:"required"`
	PartnerOPMobileNetworkCodes  MobileNetworkCodes `json:"partnerOPMobileNetworkCodes"`
	PartnerOPFixedNetworkCodes   []string           `json:"partnerOPFixedNetworkCodes"`
	FederationContextId          string             `json:"federationContextId" binding:"required"`
	EdgeDiscoveryServiceEndpoint ServiceEndpoint    `json:"edgeDiscoveryServiceEndpoint" binding:"required"`
	OfferedAvailabilityZones     []string           `json:"offeredAvailabilityZones"`
	PlatformCaps                 []string           `json:"platformCaps" binding:"required"`
	FederationRenewalDate        time.Time          `json:"federationRenewalDate"`
	FederationExpiryDate         time.Time          `json:"federationExpiryDate"`
	DefaultSubscriptions         []string           `json:"defaultSubscriptions"`
	SubscriptionsPeriodicity     string             `json:"subscriptionsPeriodicity"`
	OpsInfoExposureEndpoint      string             `json:"opsInfoExposureEndpoint"`
}
