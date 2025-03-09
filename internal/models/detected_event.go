package models

import (
	"time"
)

type DetectedEvent struct {
	ZoneId string `json:"zoneId"`

	EventId string `json:"eventId"`

	StartTime *time.Time `json:"startTime"`

	EndTime *time.Time `json:"endTime"`

	NumOccurance int32 `json:"numOccurance"`
}
