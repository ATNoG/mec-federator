package dto

type AccessTokenRequestData struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
