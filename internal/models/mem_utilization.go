package models

type MemUtilization struct {
	NoOfSamples string `json:"noOfSamples"`

	AverageUtilization *UtilizationValue `json:"averageUtilization"`

	MaxUtilization *UtilizationValue `json:"maxUtilization"`

	MinUtilization *UtilizationValue `json:"minUtilization"`

	EffectiveUtilization *UtilizationValue `json:"effectiveUtilization,omitempty"`
}
