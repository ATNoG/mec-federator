package models

type HttpMethods string

// List of HttpMethods
const (
	POST   HttpMethods = "POST"
	PUT    HttpMethods = "PUT"
	PATCH  HttpMethods = "PATCH"
	DELETE HttpMethods = "DELETE"
	GET    HttpMethods = "GET"
)
