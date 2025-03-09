package models

type Status string

const (
	FAILED            Status = "FAILED"
	TEMPORARY_FAILURE Status = "TEMPORARY_FAILURE"
	AVAILABLE         Status = "AVAILABLE"
	LOCKED            Status = "LOCKED"
	NOT_AVAILABLE     Status = "NOT_AVAILABLE"
)
