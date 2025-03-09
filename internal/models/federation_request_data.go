package models

import (
	"time"
)

type FederationRequestData struct {
	OrigOPFederationId string `json:"origOPFederationId,omitempty"`

	OrigOPCountryCode string `json:"origOPCountryCode,omitempty"`

	OrigOPMobileNetworkCodes *MobileNetworkIds `json:"origOPMobileNetworkCodes,omitempty"`

	OrigOPFixedNetworkCodes *[]string `json:"origOPFixedNetworkCodes,omitempty"`
	// Time zone info of the federation initiated by the originating OP
	InitialDate time.Time `json:"initialDate"`

	PartnerStatusLink string `json:"partnerStatusLink"`
}
