package models

type ServiceEndpoint struct {
	Port          int      `json:"port" binding:"required"`
	Fqdn          string   `json:"fqdn"`
	Ipv4Addresses []string `json:"ipv4Addresses"`
	Ipv6Addresses []string `json:"ipv6Addresses"`
}
