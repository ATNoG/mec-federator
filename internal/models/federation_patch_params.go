package models

type FederationPatchParams struct {
	// Must be either "MOBILE_NETWORK_CODES" or "FIXED_NETWORK_CODES"
	ObjectType string `json:"objectType"`
	// Must be either "ADD_CODES", "REMOVE_CODES" or "UPDATE_CODES"
	OperationType string `json:"operationType"`

	ModificationDate string `json:"modificationDate"`

	AddMobileNetworkIds *MobileNetworkIds `json:"addMobileNetworkIds,omitempty"`

	RemoveMobileNetworkIds *MobileNetworkIds `json:"removeMobileNetworkIds,omitempty"`

	AddFixedNetworkIds *[]string `json:"addFixedNetworkIds,omitempty"`

	RemoveFixedNetworkIds *[]string `json:"removeFixedNetworkIds,omitempty"`
}
