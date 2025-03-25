package models

type MobileNetworkIds struct {
	Mcc string `json:"mcc,omitempty" bson:"mcc,omitempty"`

	Mncs []string `json:"mncs,omitempty" bson:"mncs,omitempty"`
}
