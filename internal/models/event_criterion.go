package models

type EventCriterion struct {
	ResUsageType     *ResourceType `json:"resUsageType"`
	TriggerCondition string        `json:"triggerCondition"`

	ThresholdVal *ThresholdVal `json:"thresholdVal"`
	NumOccurance int32         `json:"numOccurance"`

	MonitorDuration *PeriodicityInterval `json:"monitorDuration"`
}
