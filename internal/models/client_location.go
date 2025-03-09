package models

type ClientLocation struct {
	// Latitude, Longitude as decimal fraction up to 4 digit precision
	GeoLocation string `json:"geo_location,omitempty"`
	// Information about the 4G/5G Cell ids where the client is currently served.
	RadLocation []interface{} `json:"rad_location,omitempty"`
}
