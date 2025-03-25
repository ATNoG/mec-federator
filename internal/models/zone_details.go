package models

type ZoneDetails struct {
	ZoneId string `json:"zoneId" bson:"zoneId"`

	Geolocation string `json:"geolocation" bson:"geolocation"`
	// Details about cities or state covered by the edge. Details about the type of locality for eg rural, urban, industrial etc. This information is defined in human readable form.
	GeographyDetails string `json:"geographyDetails" bson:"geographyDetails"`
}
