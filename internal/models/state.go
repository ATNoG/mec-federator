package models

type State struct {
	// Defines the alarm state during its life cycle (raised | updated | cleared).
	AlarmState string `json:"alarmState" bson:"alarmState"`
}
