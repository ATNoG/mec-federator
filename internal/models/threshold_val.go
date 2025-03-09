package models

type ThresholdVal struct {
	Value string `json:"value"`
	// The unit of resources measurement e.g. number of cores, mega bits per seconds etc.
	Unit string `json:"unit"`
}
