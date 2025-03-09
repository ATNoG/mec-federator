package models

type PatchableParams string

// List of PatchableParams
const (
	PERCEIVED_SEVERITY PatchableParams = "/perceivedSeverity"
	PROBABLE_CAUSE     PatchableParams = "/probableCause"
	ALARMED_OBJECT     PatchableParams = "/alarmedObject"
	SOURCE_SYSTEM_ID   PatchableParams = "/sourceSystemId"
	STATE              PatchableParams = "/state"
	AFFECTED_SERVICE   PatchableParams = "/affectedService"
	ALARM_DETAILS      PatchableParams = "/alarmDetails"
	SPECIFIC_PROBLEM   PatchableParams = "/specificProblem"
	SERVICE_AFFECTING  PatchableParams = "/serviceAffecting"
)
