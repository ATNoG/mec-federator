package models

type HttpMethods string

const (
	POST   HttpMethods = "POST"
	PUT    HttpMethods = "PUT"
	PATCH  HttpMethods = "PATCH"
	DELETE HttpMethods = "DELETE"
	GET    HttpMethods = "GET"
)
