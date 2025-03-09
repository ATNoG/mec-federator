package models

type ServiceEndpoint struct {
	Port int32 `json:"port"`

	Fqdn string `json:"fqdn,omitempty"`

	Ipv4Addresses []string `json:"ipv4Addresses,omitempty"`

	Ipv6Addresses []string `json:"ipv6Addresses,omitempty"`
}
