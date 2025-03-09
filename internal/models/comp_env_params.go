package models

// Environment variables are key value pairs that should be injected when component in instantiated
type CompEnvParams struct {
	// Name of environment variable
	EnvVarName string `json:"envVarName"`

	EnvValueType string `json:"envValueType"`
	// Value to be assigned to environment variable
	EnvVarValue string `json:"envVarValue,omitempty"`
	// Full path of parameter from componentSpec that should be used to generate the environment value. Eg. networkResourceProfile[1]. interfaceId.
	EnvVarSrc string `json:"envVarSrc,omitempty"`
}
