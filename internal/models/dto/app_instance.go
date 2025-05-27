package dto

type InstantiateApplicationRequest struct {
	TransactionId       string `json:"transactionId"`
	AppId               string `json:"appId"`
	AppProviderId       string `json:"appProviderId"`
	AppVersion          string `json:"appVersion"`
	ZoneInfo            string `json:"zoneInfo"`
	AppInstCallbackLink string `json:"appInstCallbackLink"`
}
