package models

type ResourceReservationDuration struct {
	NumOfDays   int32 `json:"numOfDays,omitempty"`
	NumOfMonths int32 `json:"numOfMonths,omitempty"`
	NumOfYears  int32 `json:"numOfYears,omitempty"`
}
