package models

type FederationInitiateRequestData struct {
	FederationEndpoint string `json:"federationEndpoint" binding:"required"`
	AuthEndpoint       string `json:"authEndpoint" binding:"required"`
}
