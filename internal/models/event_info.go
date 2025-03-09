package models

type EventInfo struct {
	EventId string `json:"eventId"`

	EventCriterion *EventCriterion `json:"eventCriterion"`
}
