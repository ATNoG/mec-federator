package models

import "time"

type AccessToken struct {
	ClientId    string    `json:"clientId" bson:"clientId"`
	AccessToken string    `json:"accessToken" bson:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt" bson:"expiresAt"`
}
