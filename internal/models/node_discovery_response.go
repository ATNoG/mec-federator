package models

type NodeDiscoveryResponse struct {
	EdgeNodes *[]interface{} `json:"edgeNodes"`

	DiscoveredAppInsts *[]interface{} `json:"discoveredAppInsts"`
}
