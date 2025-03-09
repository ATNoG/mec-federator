package models

type ExpiryInterval struct {
	NumHours int32 `json:"numHours"`
	NumMins  int32 `json:"numMins"`
	NumSecs  int32 `json:"numSecs"`
}
