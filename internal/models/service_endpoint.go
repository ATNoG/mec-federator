package models

type ServiceEndpoint struct {
	Port int32 `json:"port" bson:"port"`

	Fqdn string `json:"fqdn,omitempty" bson:"fqdn,omitempty"`

	Ipv4Addresses []string `json:"ipv4Addresses,omitempty" bson:"ipv4Addresses,omitempty"`

	Ipv6Addresses []string `json:"ipv6Addresses,omitempty" bson:"ipv6Addresses,omitempty"`
}
