package models

type UpdatedParam struct {
	AlarmUpdateOps *AlarmUpdateOps `json:"alarmUpdateOps"`

	PatchableParam *PatchableParams `json:"patchableParam"`
	PatchValue     string           `json:"patchValue"`
}
