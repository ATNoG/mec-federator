package models

type CpuUtilization struct {
	CpuType *MonitoringSubsType `json:"cpuType"`
	// Number of samples used for calculating metrics.
	NoOfSamples string `json:"noOfSamples"`

	AverageUtilization *UtilizationValue `json:"averageUtilization"`

	MaxUtilization *UtilizationValue `json:"maxUtilization"`

	MinUtilization *UtilizationValue `json:"minUtilization"`

	EffectiveUtilization *UtilizationValue `json:"effectiveUtilization"`
}
