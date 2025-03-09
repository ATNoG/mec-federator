package models

type PersistentVolumeDetails struct {
	// size of the volume given by user (10GB, 20GB, 50 GB or 100GB)
	VolumeSize string `json:"volumeSize"`
	// Defines the mount path of the volume
	VolumeMountPath string `json:"volumeMountPath"`
	// Human readable name for the volume
	VolumeName string `json:"volumeName"`
	// It indicates the ephemeral storage on the node and contents are not preserved if containers restarts
	EphemeralType bool `json:"ephemeralType,omitempty"`
	// Values are RW (read/write) and RO (read-only)l
	AccessMode string `json:"accessMode,omitempty"`
	// Exclusive or Shared. If shared, then in case of multiple containers same volume will be shared across the containers.
	SharingPolicy string `json:"sharingPolicy,omitempty"`
}
