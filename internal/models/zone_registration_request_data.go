package models

type ZoneRegistrationRequestData struct {
	FederationContextId string `json:"federationContextId"`

	AcceptedAvailabilityZones []string `json:"acceptedAvailabilityZones"`

	AvailZoneNotifLink string `json:"availZoneNotifLink"`
}
