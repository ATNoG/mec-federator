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
	KduId  string `json:"kduId"`
	NsId   string `json:"nsId"`
	Node   string `json:"node"`
	NodeId string `json:"nodeId"`
}

// request to disable a KDU of an application instance
type DisableAppInstanceKDURequest struct {
	NsId  string `json:"nsId"`
	KduId string `json:"kduId"`
}
