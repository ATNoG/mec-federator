package models

type FlavourMetrics struct {
	NoOfSamples string `json:"noOfSamples"`

	FlavourId string `json:"flavourId"`

	AverageUtilization *UtilizationValue `json:"averageUtilization"`

	AverageThroughput *UtilizationValue `json:"averageThroughput,omitempty"`

	MaxUtilization *UtilizationValue `json:"maxUtilization"`

	MinUtilization *UtilizationValue `json:"minUtilization"`
}
