package models

type AlarmedObject struct {
	AlarmId *AlarmIdentifier `json:"alarmId"`

	Href string `json:"href"`
}
