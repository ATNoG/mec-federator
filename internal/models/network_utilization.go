package models

type NetworkUtilization struct {
	NoOfSamples string `json:"noOfSamples"`

	IngressUsage *UtilizationValue `json:"ingressUsage"`

	EgressUsage *UtilizationValue `json:"egressUsage"`

	AverageThroughput *UtilizationValue `json:"averageThroughput"`

	MaxThroughput *UtilizationValue `json:"maxThroughput"`

	MinThroughput *UtilizationValue `json:"minThroughput"`
}
