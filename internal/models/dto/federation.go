package dto

import "time"

type FederationInitiateRequestData struct {
	FederationEndpoint string `json:"federationEndpoint" binding:"required"`
	AuthEndpoint       string `json:"authEndpoint" binding:"required"`
}

type FederationRenewalResponseData struct {
	FederationContextId   string    `json:"federationContextId"`
	FederationExpiryDate  time.Time `json:"federationExpiryDate"`
	FederationRenewalDate time.Time `json:"federationRenewalDate"`
}
