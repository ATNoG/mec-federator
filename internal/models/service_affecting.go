package models

// ServiceAffecting : Specific information related to the alarm
type ServiceAffecting string

// List of ServiceAffecting
const (
	YES ServiceAffecting = "YES"
	NO  ServiceAffecting = "NO"
)
