package models

type PeriodicNotifConfig struct {
	Periodicity *PeriodicityInterval `json:"periodicity,omitempty"`

	NotificationListner string `json:"notificationListner,omitempty"`
}
