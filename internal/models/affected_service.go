package models

type AffectedService struct {
	// Defines the affected services e.g., edge discovery, application services, API services etc at source
	AffectedService []string `json:"affectedService"`
}
