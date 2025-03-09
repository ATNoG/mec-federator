package models

type LocationRetrievalCap struct {
	CapabilityId *CapabilityId `json:"capabilityId"`
	// The enumerated list of UE location accuracy that an OP can determine via SBI-NR.
	LocationType string `json:"locationType"`
	// The enumerated list of type of network location of an UE that an OP can determine via SBI-NR.
	LocationAccuracy string `json:"locationAccuracy,omitempty"`
}
