package models

type HttpResources struct {
	Href        string        `json:"href"`
	HttpMethods []HttpMethods `json:"httpMethods"`
}
