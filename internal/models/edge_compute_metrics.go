package models

import (
	"time"
)

type EdgeComputeMetrics struct {
	ZoneId string `json:"zoneId"`

	StartTime *time.Time `json:"startTime"`

	EndTime *time.Time `json:"endTime"`

	CpuUtil *CpuUtilization `json:"cpuUtil"`

	MemUtil *MemUtilization `json:"memUtil"`

	DiskUtil *DiskUtilization `json:"diskUtil"`

	NetworkUtil *NetworkUtilization `json:"networkUtil"`

	FlavourUtil *[]FlavourMetrics `json:"flavourUtil"`
}
