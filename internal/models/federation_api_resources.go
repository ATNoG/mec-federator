package models

type FederationApiResources struct {
	Name          *FederationApiNames `json:"name"`
	ApiOperations []HttpResources     `json:"apiOperations"`
}
