package models

type MonitoringSubsType string

const (
	EDGE_RESOURCE MonitoringSubsType = "edge_resource"
	APP_RESOURCE  MonitoringSubsType = "app_resource"
	ALARM         MonitoringSubsType = "alarm"
	ALL           MonitoringSubsType = "all"
)
