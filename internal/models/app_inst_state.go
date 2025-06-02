package models

// InstanceState : Running status of the application instance.
type AppInstState string

// List of InstanceState
const (
	INSTANCE_PENDING     AppInstState = "PENDING"
	INSTANCE_READY       AppInstState = "READY"
	INSTANCE_FAILED      AppInstState = "FAILED"
	INSTANCE_TERMINATING AppInstState = "TERMINATING"
)
