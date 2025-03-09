package models

type EventSubscriptionInfo struct {
	ResUsageType *ResourceType `json:"resUsageType"`

	Periodicity *PeriodicityInterval `json:"periodicity"`

	SubscriptionId string `json:"subscriptionId"`
}
