package nbi

import "time"

// FederationManagement

type FederationInitiateRequest struct {
	FederationEndpoint string `json:"federationEndpoint" binding:"required"`
	AuthEndpoint       string `json:"authEndpoint" binding:"required"`
}

type FederationRenewalResponse struct {
	FederationContextId   string    `json:"federationContextId"`
	FederationExpiryDate  time.Time `json:"federationExpiryDate"`
	FederationRenewalDate time.Time `json:"federationRenewalDate"`
}

// Artefact Management

