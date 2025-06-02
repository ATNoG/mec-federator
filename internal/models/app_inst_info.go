package models

type AppInstInfo struct {
	AppInstState    AppInstState `json:"appInstState"`
	AccessPointInfo []string     `json:"accessPointInfo"`
}
