package models

type ResourceType string

const (
	CPU     ResourceType = "CPU"
	MEMORY  ResourceType = "MEMORY"
	DISK    ResourceType = "DISK"
	NETWORK ResourceType = "NETWORK"
	FLAVOUR ResourceType = "FLAVOUR"
)
