package models

type OsType struct {
	Architecture string `json:"architecture"`

	Distribution string `json:"distribution"`

	Version string `json:"version"`

	License string `json:"license"`
}
