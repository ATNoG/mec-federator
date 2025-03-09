package models

// MonitoringSubsType : Denotes types of edge resources, faults and events at partner OP to be reported to Originating OP.
type MonitoringSubsType string

// List of monitoringSubsType
const (
	EDGE_RESOURCE MonitoringSubsType = "edge_resource"
	APP_RESOURCE  MonitoringSubsType = "app_resource"
	ALARM         MonitoringSubsType = "alarm"
	ALL           MonitoringSubsType = "all"
)
