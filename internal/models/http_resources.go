package models

type HttpResources struct {
	Href string `json:"href"`
	// List of HTTP Methods supported for the given API category
	HttpMethods []HttpMethods `json:"httpMethods"`
}
