package models

type CompEnvParams struct {
	EnvVarName string `json:"envVarName"`

	EnvValueType string `json:"envValueType"`
	EnvVarValue  string `json:"envVarValue,omitempty"`
	EnvVarSrc    string `json:"envVarSrc,omitempty"`
}
