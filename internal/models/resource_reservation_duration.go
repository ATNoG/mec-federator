package models

// Time period for which resources are to be reserved starting from now
type ResourceReservationDuration struct {
	// Number of days to be reserved
	NumOfDays int32 `json:"numOfDays,omitempty"`
	// Number of months to be reserved
	NumOfMonths int32 `json:"numOfMonths,omitempty"`
	// Number of years to be reserved
	NumOfYears int32 `json:"numOfYears,omitempty"`
}
