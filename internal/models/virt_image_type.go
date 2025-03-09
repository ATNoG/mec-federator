package models

type VirtImageType string

const (
	QCOW2  VirtImageType = "QCOW2"
	DOCKER VirtImageType = "DOCKER"
	OVA    VirtImageType = "OVA"
)
