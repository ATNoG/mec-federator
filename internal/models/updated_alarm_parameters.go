package models

type UpdatedAlarmParameters struct {
	AlarmId      *AlarmIdentifier `json:"alarmId"`
	UpdateParams []UpdatedParam   `json:"updateParams"`
}
