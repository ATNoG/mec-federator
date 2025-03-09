package models

type OnboardStatusInfo string

const (
	ONBOARD_PENDING    OnboardStatusInfo = "PENDING"
	ONBOARD_ONBOARDED  OnboardStatusInfo = "ONBOARDED"
	ONBOARD_DEBOARDING OnboardStatusInfo = "DEBOARDING"
	ONBOARD_REMOVED    OnboardStatusInfo = "REMOVED"
	ONBOARD_FAILED     OnboardStatusInfo = "FAILED"
)
