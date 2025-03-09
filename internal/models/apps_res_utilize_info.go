package models

type AppsResUtilizeInfo struct {
	AppMetrics []AppsResUtilizeMetrics `json:"appMetrics"`

	FederationContextId string `json:"federationContextId"`
	SequenceNum         int32  `json:"sequenceNum"`
}
