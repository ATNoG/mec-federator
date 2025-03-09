package models

type FederationSupportedApis struct {
	FederationBaseAPI *FederationApiResources `json:"federationBaseAPI"`

	AvailabilityZoneAPI *FederationApiResources `json:"availabilityZoneAPI"`

	EdgeApplicationAPI *FederationApiResources `json:"edgeApplicationAPI"`

	ArtefactAPI *FederationApiResources `json:"artefactAPI"`

	FileAPI *FederationApiResources `json:"fileAPI"`

	ServiceAPIFederation *FederationApiResources `json:"serviceAPIFederation,omitempty"`

	ResourceMonitoringAPI *FederationApiResources `json:"resourceMonitoringAPI,omitempty"`

	FaultManagementAPI *FederationApiResources `json:"faultManagementAPI,omitempty"`

	EventManagementAPI *FederationApiResources `json:"eventManagementAPI,omitempty"`
}
