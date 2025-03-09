package models

type LocationRetrievalCap struct {
	CapabilityId     *CapabilityId `json:"capabilityId"`
	LocationType     string        `json:"locationType"`
	LocationAccuracy string        `json:"locationAccuracy,omitempty"`
}
