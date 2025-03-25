package models

import "time"

type FederationRenewalResponseData struct {
	FederationContextId string `json:"federationContextId"`

	FederationExpiryDate time.Time `json:"federationExpiryDate"`

	FederationRenewalDate time.Time `json:"federationRenewalDate"`
}
