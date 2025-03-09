package models

// Parameters corresponding to the performance constraints, tenancy details etc.
type AppQoSProfile struct {
	LatencyConstraints *LatencyConstraints `json:"latencyConstraints"`

	BandwidthRequired int32 `json:"bandwidthRequired,omitempty"`

	MultiUserClients *MultiUserClients `json:"multiUserClients,omitempty"`

	NoOfUsersPerAppInst int32 `json:"noOfUsersPerAppInst,omitempty"`

	AppProvisioning bool `json:"appProvisioning,omitempty"`
}
