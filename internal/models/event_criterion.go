package models

type EventCriterion struct {
	ResUsageType *ResourceType `json:"resUsageType"`
	// The condition evaluation operator to compare threashold value of a resource for event detection.
	TriggerCondition string `json:"triggerCondition"`

	ThresholdVal *ThresholdVal `json:"thresholdVal"`
	// Number of times the trigger condition is detected
	NumOccurance int32 `json:"numOccurance"`

	MonitorDuration *PeriodicityInterval `json:"monitorDuration"`
}
