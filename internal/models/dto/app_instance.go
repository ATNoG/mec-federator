package dto

// request to instantiate an application
type InstantiateApplicationRequest struct {
	TransactionId       string `json:"transactionId"`
	AppId               string `json:"appId"`
	AppProviderId       string `json:"appProviderId"`
	AppVersion          string `json:"appVersion"`
	ZoneInfo            string `json:"zoneInfo"`
	AppInstCallbackLink string `json:"appInstCallbackLink"`
	Config              string `json:"config"`
}

// request to enable a KDU of an application instance
type EnableAppInstanceKDURequest struct {
	KduId string `json:"kduId"`
	NsId  string `json:"nsId"`
	Node  string `json:"node"`
}

// request to disable a KDU of an application instance
type DisableAppInstanceKDURequest struct {
	NsId  string `json:"nsId"`
	KduId string `json:"kduId"`
}

// request to migrate an application instance to a specific node
type AppInstanceNodeMigrateRequest struct {
	NsId  string `json:"nsId"`
	VnfId string `json:"vnfId"`
	KduId string `json:"kduId"`
	Node  string `json:"node"`
}
