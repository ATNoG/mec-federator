package models

type PersistentVolumeDetails struct {
	VolumeSize      string `json:"volumeSize"`
	VolumeMountPath string `json:"volumeMountPath"`
	VolumeName      string `json:"volumeName"`
	EphemeralType   bool   `json:"ephemeralType,omitempty"`
	AccessMode      string `json:"accessMode,omitempty"`
	SharingPolicy   string `json:"sharingPolicy,omitempty"`
}
