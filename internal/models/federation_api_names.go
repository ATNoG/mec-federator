package models

type FederationApiNames string

// List of FederationAPINames
const (
	FEDERATION FederationApiNames = "FEDERATION"
	AVAILZONE  FederationApiNames = "AVAILZONE"
	ARTEFACT   FederationApiNames = "ARTEFACT"
	FILE       FederationApiNames = "FILE"
	SVSAPEFED  FederationApiNames = "SVSAPEFED"
	RESMONITOR FederationApiNames = "RESMONITOR"
	EVENTMGMT  FederationApiNames = "EVENTMGMT"
	FAULTMGMT  FederationApiNames = "FAULTMGMT"
)
