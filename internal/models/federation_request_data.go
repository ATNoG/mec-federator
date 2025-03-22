package models

import (
	"time"
)

type FederationRequestData struct {
	OrigOPFederationId string `json:"origOPFederationId,omitempty" bson:"origOPFederationId,omitempty"`

	OrigOPCountryCode string `json:"origOPCountryCode,omitempty" bson:"origOPCountryCode,omitempty"`

	OrigOPMobileNetworkCodes *MobileNetworkIds `json:"origOPMobileNetworkCodes,omitempty" bson:"origOPMobileNetworkCodes,omitempty"`

	OrigOPFixedNetworkCodes *[]string `json:"origOPFixedNetworkCodes,omitempty" bson:"origOPFixedNetworkCodes,omitempty"`
	// Time zone info of the federation initiated by the originating OP
	InitialDate time.Time `json:"initialDate" bson:"initialDate" binding:"required"`

	PartnerStatusLink string `json:"partnerStatusLink" bson:"partnerStatusLink" binding:"required"`

	AccessToken AccessToken `json:"accessToken" bson:"accessToken" binding:"required"`
}
