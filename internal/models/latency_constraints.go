package models

type LatencyConstraints string

const (
	NONE     LatencyConstraints = "NONE"
	LOW      LatencyConstraints = "LOW"
	ULTRALOW LatencyConstraints = "ULTRALOW"
)
