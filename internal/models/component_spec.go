package models

// Details about compute, networking and storage requirements for each component of the application. App provider should define all information needed to instantiate the component. If artefact is being defined at component level this section should have information just about the component. In case the artefact is being defined at application level the section should provide details about all the components.
type ComponentSpec struct {
	// Must be a valid RFC 1035 label name.  Component name must be unique with an application
	ComponentName string `json:"componentName"`
	// List of all images associated with the component. Images are specified using the file identifiers. Partner OP provides these images using file upload api.
	Images []string `json:"images"`
	// Number of component instances to be launched.
	NumOfInstances int32 `json:"numOfInstances"`
	// How the platform shall handle component failure
	RestartPolicy string `json:"restartPolicy"`

	CommandLineParams *CommandLineParams `json:"commandLineParams,omitempty"`
	// Each application component exposes some ports either for external users or for inter component communication. Application provider is required to specify which ports are to be exposed and the type of traffic that will flow through these ports.
	ExposedInterfaces []InterfaceDetails `json:"exposedInterfaces,omitempty"`

	ComputeResourceProfile *ComputeResourceInfo `json:"computeResourceProfile"`

	CompEnvParams []CompEnvParams `json:"compEnvParams,omitempty"`

	DeploymentConfig *DeploymentConfig `json:"deploymentConfig,omitempty"`
	// The ephemeral volume a container process may need to temporary store internal data
	PersistentVolumes []PersistentVolumeDetails `json:"persistentVolumes,omitempty"`
}
