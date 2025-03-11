package models

type Federation struct {
	OriginOP FederationRequestData `json:"originOP" bson:"originOP" binding:"required"`

	PartnerOP FederationResponseData `json:"partnerOP" bson:"partnerOP" binding:"required"`

	HealthInfo FederationHealthInfo `json:"federationHealthInfo,omitempty" bson:"federationHealthInfo"`

	IsEstablished bool `json:"status,omitempty" bson:"status"`

	IsOriginOP bool `json:"-" bson:"isOriginOP"`
}
