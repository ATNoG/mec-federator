package models

type InstanceState string

const (
	INSTANCE_PENDING     InstanceState = "PENDING"
	INSTANCE_READY       InstanceState = "READY"
	INSTANCE_FAILED      InstanceState = "FAILED"
	INSTANCE_TERMINATING InstanceState = "TERMINATING"
)
