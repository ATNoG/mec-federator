package models

type AlarmObjectInfo struct {
	AlarmType *AlarmType `json:"alarmType"`

	AlarmId *AlarmIdentifier `json:"alarmId"`

	PerceivedSeverity *PerceivedSeverity `json:"perceivedSeverity"`

	ProbableCause *ProbableCause `json:"probableCause"`

	AlarmedObject *AlarmedObject `json:"alarmedObject"`

	SourceSystemId *SourceSystemId `json:"sourceSystemId"`

	State *State `json:"state"`

	AlarmRaisedTime *AlarmRaisedTime `json:"alarmRaisedTime"`

	AffectedService *AffectedService `json:"affectedService,omitempty"`

	AlarmDetails *AlarmDetails `json:"alarmDetails,omitempty"`

	SpecificProblem *SpecificProblem `json:"specificProblem,omitempty"`

	ServiceAffecting *ServiceAffecting `json:"serviceAffecting,omitempty"`
}
