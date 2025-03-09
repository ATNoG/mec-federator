package models

type EdgeResUtilizeMetrics struct {
	// List of edge cloud resource metrics per zone
	EdgeMetrics []EdgeComputeMetrics `json:"edgeMetrics"`

	FederationContextId string `json:"federationContextId"`
	// Monotonically increasing counter for sequencing resource monitoring reports
	SequenceNum int32 `json:"sequenceNum"`
}
