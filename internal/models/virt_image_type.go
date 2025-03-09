package models

// VirtImageType : Indicate if the file is Container image or VM image (QCOW2, OVA)
type VirtImageType string

// List of VirtImageType
const (
	QCOW2  VirtImageType = "QCOW2"
	DOCKER VirtImageType = "DOCKER"
	OVA    VirtImageType = "OVA"
)
