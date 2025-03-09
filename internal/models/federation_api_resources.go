package models

type FederationApiResources struct {
	Name *FederationApiNames `json:"name"`
	// List of HTTP Methods supported for the given API category
	ApiOperations []HttpResources `json:"apiOperations"`
}
