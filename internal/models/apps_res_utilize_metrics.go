package models

import (
	"time"
)

type AppsResUtilizeMetrics struct {
	ZoneId string `json:"zoneId"`

	StartTime *time.Time `json:"startTime"`

	EndTime *time.Time `json:"endTime"`

	AppZoneMetrics *[]AppAggrResUtil `json:"appZoneMetrics"`
}
