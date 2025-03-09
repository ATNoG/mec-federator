package models

// Configuration used when deploying a component. May override other ComponentSpec parameters related to deployment like restart policy, command line parameters, environment variables, etc.
type DeploymentConfig struct {
	// Config type.
	ConfigType string `json:"configType"`
	// Contents of the configuration.
	Contents string `json:"contents"`
}
