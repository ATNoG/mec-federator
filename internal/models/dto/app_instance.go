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
	KDUId string `json:"kduId"`
	Node  string `json:"node"`
}

// request to disable a KDU of an application instance
type DisableAppInstanceKDURequest struct {
	KDUId string `json:"kduId"`
}
