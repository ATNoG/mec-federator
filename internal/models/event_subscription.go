package models

type EventSubscription struct {
	ResUsageType *ResourceType `json:"resUsageType"`

	Periodicity *PeriodicityInterval `json:"periodicity"`

	EventListner string `json:"eventListner"`
}
