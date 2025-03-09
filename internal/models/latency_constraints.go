package models

// LatencyConstraints : Latency requirements for the application.Allowed values (non-standardized) are none, low and ultra-low. Ultra-Low may corresponds to range 15 - 30 msec, Low correspond to range 30 - 50 msec. None means 51 and above
type LatencyConstraints string

// List of LatencyConstraints
const (
	NONE     LatencyConstraints = "NONE"
	LOW      LatencyConstraints = "LOW"
	ULTRALOW LatencyConstraints = "ULTRALOW"
)
