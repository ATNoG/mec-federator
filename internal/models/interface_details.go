package models

type InterfaceDetails struct {
	InterfaceId    string `json:"interfaceId"`
	CommProtocol   string `json:"commProtocol"`
	CommPort       int32  `json:"commPort"`
	VisibilityType string `json:"visibilityType"`
	Network        string `json:"network,omitempty"`
	InterfaceName  string `json:"InterfaceName,omitempty"`
}
