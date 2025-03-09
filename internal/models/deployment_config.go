package models

type DeploymentConfig struct {
	ConfigType string `json:"configType"`
	Contents   string `json:"contents"`
}
