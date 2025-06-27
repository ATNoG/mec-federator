package models

type ZoneRegistrationRequestData struct {
	AcceptedAvailabilityZones []string `json:"acceptedAvailabilityZones"`

	AvailZoneNotifLink string `json:"availZoneNotifLink"`
}
