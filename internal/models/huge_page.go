package models

type HugePage struct {
	// Size of hugepage
	PageSize string `json:"pageSize"`
	// Total number of huge pages
	Number int32 `json:"number"`
}
