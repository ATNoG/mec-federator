package models

import (
	"time"
)

type AlarmRaisedTime struct {
	// Defines the alarm raised time at source
	AlarmRaisedTime time.Time `json:"alarmRaisedTime"`
}
