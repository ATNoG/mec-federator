package models

type PeriodicityInterval struct {
	// Number of Hours for Expiry (0-23)
	NumHours int32 `json:"numHours"`
	// Number of Minutes for Expiry (0-59)
	NumMins int32 `json:"numMins"`
}
