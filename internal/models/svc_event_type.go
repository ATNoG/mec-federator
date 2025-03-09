package models

type SvcEventType string

// List of svcEventType
const (
	EVENT_TIMEREXPIRY SvcEventType = "evt_timerexpiry"
	EVENT_NETWORK     SvcEventType = "evt_network"
	EVENT_DELETE      SvcEventType = "evt_delete"
)
