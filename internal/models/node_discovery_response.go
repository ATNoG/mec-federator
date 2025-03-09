package models

// Candidate availability zones and details of already running instances of the given application
type NodeDiscoveryResponse struct {
	EdgeNodes *[]interface{} `json:"edgeNodes"`

	DiscoveredAppInsts *[]interface{} `json:"discoveredAppInsts"`
}
