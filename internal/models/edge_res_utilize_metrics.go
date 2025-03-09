package models

type EdgeResUtilizeMetrics struct {
	EdgeMetrics []EdgeComputeMetrics `json:"edgeMetrics"`

	FederationContextId string `json:"federationContextId"`
	SequenceNum         int32  `json:"sequenceNum"`
}
