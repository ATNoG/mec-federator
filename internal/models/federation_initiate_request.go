package models

type FederationInitiateRequest struct {
	FederationEndpoint string `json:"federationEndpoint" binding:"required"`
	AuthEndpoint       string `json:"authEndpoint" binding:"required"`
}
