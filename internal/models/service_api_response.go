package models

type ServiceApiResponse struct {
	CustomerID string `json:"customerID"`

	TargetUserContext *TargetUserContext `json:"targetUserContext"`

	ApiResponse string `json:"apiResponse"`

	TxnIdentifier string `json:"txnIdentifier"`
}
