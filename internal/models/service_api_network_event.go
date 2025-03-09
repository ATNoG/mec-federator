package models

type ServiceApiNetworkEvent struct {
	ConnectID string `json:"connectID"`

	CustomerID string `json:"customerID"`

	EventType *SvcEventType `json:"EventType"`

	ServiceAPIEventDef *ServiceApiEventDef `json:"serviceAPIEventDef,omitempty"`

	ExpiryDuration *ExpiryInterval `json:"expiryDuration,omitempty"`
}
