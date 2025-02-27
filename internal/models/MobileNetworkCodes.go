package models

type MobileNetworkCodes struct {
	Mcc  string   `json:"mcc"`
	Mncs []string `json:"mncs"`
}
