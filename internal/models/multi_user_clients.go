package models

// MultiUserClients : Single user type application are designed to serve just one client. Multi user type application is designed to serve multiple clients
type MultiUserClients string

// List of MultiUserClients
const (
	SINGLE_USER MultiUserClients = "APP_TYPE_SINGLE_USER"
	MULTI_USER  MultiUserClients = "APP_TYPE_MULTI_USER"
)
