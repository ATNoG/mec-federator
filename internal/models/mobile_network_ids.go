package models

type MobileNetworkIds struct {
	Mcc string `json:"mcc,omitempty"`

	Mncs []string `json:"mncs,omitempty"`
}
