package models

import "time"

type AccessToken struct {
	AccessToken string    `json:"accessToken" bson:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt" bson:"expiresAt"`
}
