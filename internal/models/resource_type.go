package models

// ResourceType : Indicate the type of resource
type ResourceType string

// List of resourceType
const (
	CPU     ResourceType = "CPU"
	MEMORY  ResourceType = "MEMORY"
	DISK    ResourceType = "DISK"
	NETWORK ResourceType = "NETWORK"
	FLAVOUR ResourceType = "FLAVOUR"
)
