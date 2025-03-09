package models

// Indicate the resource measurement Unit
type UtilizationValue struct {
	ResType *ResourceType `json:"resType"`
	// Whole number that represent the value of given resource type.
	Value string `json:"value"`

	Unit string `json:"unit"`
}
