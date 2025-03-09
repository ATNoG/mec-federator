package models

import (
	"time"
)

type ResourceSubscriptionInfo struct {
	MonitoringType *MonitoringSubsType `json:"monitoringType"`

	DateAndTime    *time.Time `json:"dateAndTime"`
	SubscriptionId string     `json:"subscriptionId"`
}
