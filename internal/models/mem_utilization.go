package models

type MemUtilization struct {
	// Number of samples used for calculating metrics.
	NoOfSamples string `json:"noOfSamples"`

	AverageUtilization *UtilizationValue `json:"averageUtilization"`

	MaxUtilization *UtilizationValue `json:"maxUtilization"`

	MinUtilization *UtilizationValue `json:"minUtilization"`

	EffectiveUtilization *UtilizationValue `json:"effectiveUtilization,omitempty"`
}
