package models

import "time"

type AccessToken struct {
	Token     string    `json:"accessToken" bson:"accessToken"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
	Role      string    `json:"role" bson:"role"`
}
