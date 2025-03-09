package models

type ZoneDetails struct {
	ZoneId string `json:"zoneId"`

	Geolocation      string `json:"geolocation"`
	GeographyDetails string `json:"geographyDetails"`
}
