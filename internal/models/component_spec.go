package models

type ComponentSpec struct {
	ComponentName  string   `json:"componentName"`
	Images         []string `json:"images"`
	NumOfInstances int32    `json:"numOfInstances"`
	RestartPolicy  string   `json:"restartPolicy"`

	CommandLineParams *CommandLineParams `json:"commandLineParams,omitempty"`
	ExposedInterfaces []InterfaceDetails `json:"exposedInterfaces,omitempty"`

	ComputeResourceProfile *ComputeResourceInfo `json:"computeResourceProfile"`

	CompEnvParams []CompEnvParams `json:"compEnvParams,omitempty"`

	DeploymentConfig  *DeploymentConfig         `json:"deploymentConfig,omitempty"`
	PersistentVolumes []PersistentVolumeDetails `json:"persistentVolumes,omitempty"`
}
