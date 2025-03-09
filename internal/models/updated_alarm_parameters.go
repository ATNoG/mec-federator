package models

type UpdatedAlarmParameters struct {
	AlarmId *AlarmIdentifier `json:"alarmId"`
	// List of alarm parameters to be updated in an update operation
	UpdateParams []UpdatedParam `json:"updateParams"`
}
