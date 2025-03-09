package models

type MultiUserClients string

const (
	SINGLE_USER MultiUserClients = "APP_TYPE_SINGLE_USER"
	MULTI_USER  MultiUserClients = "APP_TYPE_MULTI_USER"
)
