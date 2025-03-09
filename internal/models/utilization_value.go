package models

type UtilizationValue struct {
	ResType *ResourceType `json:"resType"`
	Value   string        `json:"value"`

	Unit string `json:"unit"`
}
