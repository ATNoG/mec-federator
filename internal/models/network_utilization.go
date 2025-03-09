package models

type NetworkUtilization struct {
	// Number of samples used for calculating metrics.
	NoOfSamples string `json:"noOfSamples"`

	IngressUsage *UtilizationValue `json:"ingressUsage"`

	EgressUsage *UtilizationValue `json:"egressUsage"`

	AverageThroughput *UtilizationValue `json:"averageThroughput"`

	MaxThroughput *UtilizationValue `json:"maxThroughput"`

	MinThroughput *UtilizationValue `json:"minThroughput"`
}
