package models

import "time"

type FederationRequestData struct {
	OrigOPFederationId       string             `json:"origOPFederationId"`
	OrigOPCountryCode        string             `json:"origOPCountryCode"`
	OrigOPMobileNetworkCodes MobileNetworkCodes `json:"origOPMobileNetworkCodes"`
	OrigOPFixedNetworkCodes  []string           `json:"origOPFixedNetworkCodes"`
	InitialDate              time.Time          `json:"initialDate" binding:"required"`
	PartnerStatusLink        string             `json:"partnerStatusLink" binding:"required"`
}
