package models

// InstanceState : Running status of the application instance.
type InstanceState string

// List of InstanceState
const (
	INSTANCE_PENDING     InstanceState = "PENDING"
	INSTANCE_READY       InstanceState = "READY"
	INSTANCE_FAILED      InstanceState = "FAILED"
	INSTANCE_TERMINATING InstanceState = "TERMINATING"
)
