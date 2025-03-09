package models

type AppsResUtilizeInfo struct {
	// List of edge cloud resource metrics per zone
	AppMetrics []AppsResUtilizeMetrics `json:"appMetrics"`

	FederationContextId string `json:"federationContextId"`
	// Monotonically increasing counter for sequencing app monitoring reports
	SequenceNum int32 `json:"sequenceNum"`
}
