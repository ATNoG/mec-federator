package models

type ClientLocation struct {
	GeoLocation string        `json:"geo_location,omitempty"`
	RadLocation []interface{} `json:"rad_location,omitempty"`
}
