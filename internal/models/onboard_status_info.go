package models

// OnboardStatusInfo : Defines change in application status. This change could be related to application itself or an application instance status
type OnboardStatusInfo string

// List of OnboardStatusInfo
const (
	ONBOARD_PENDING    OnboardStatusInfo = "PENDING"
	ONBOARD_ONBOARDED  OnboardStatusInfo = "ONBOARDED"
	ONBOARD_DEBOARDING OnboardStatusInfo = "DEBOARDING"
	ONBOARD_REMOVED    OnboardStatusInfo = "REMOVED"
	ONBOARD_FAILED     OnboardStatusInfo = "FAILED"
)
