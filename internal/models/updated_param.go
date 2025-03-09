package models

type UpdatedParam struct {
	AlarmUpdateOps *AlarmUpdateOps `json:"alarmUpdateOps"`

	PatchableParam *PatchableParams `json:"patchableParam"`
	// Value to be replaced for the alarm parameter being updated
	PatchValue string `json:"patchValue"`
}
