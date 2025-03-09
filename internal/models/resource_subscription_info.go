package models

import (
	"time"
)

type ResourceSubscriptionInfo struct {
	MonitoringType *MonitoringSubsType `json:"monitoringType"`

	DateAndTime *time.Time `json:"dateAndTime"`
	// Partner OP managed identifier for new subscription.
	SubscriptionId string `json:"subscriptionId"`
}
