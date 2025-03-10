package models

import "time"

type AccessToken struct {
	Token             string    `json:"accessToken" bson:"accessToken"`
	ExpiresAt         time.Time `json:"expiresAt" bson:"expiresAt"`
	AuthorizationCode string    `json:"authorizationCode" bson:"authorizationCode"`
	Role              string    `json:"role" bson:"role"`
	AppProviderId     string    `json:"-" bson:"appProviderId"` // to associate token with appProviderId (this is only a valid value if the token is related to an app. provider)
}
