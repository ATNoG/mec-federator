package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

/*
 * AuthService
 *	responsible for managing access tokens
 */

// interface for authentication service
type AuthServiceInterface interface {
	SaveAccessToken(accessToken models.AccessToken) error
	QueryAccessToken(tokenStr string) (models.AccessToken, error)
}

// implementation of authentication service, requiring a mongo client
type AuthService struct {
}

// create a new authentication service, injecting given mongo client
func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) getAccessTokenCollection() *mongo.Collection {
	return config.GetMongoDatabase().Collection("access_tokens")
}

// save an access token to the database
func (s *AuthService) SaveAccessToken(accessToken models.AccessToken) error {
	collection := s.getAccessTokenCollection()
	_, err := collection.InsertOne(context.TODO(), accessToken)
	log.Printf("AuthService - Saved access token")
	return err
}

// query an access token from the database; if experired, delete it
func (s *AuthService) QueryAccessToken(tokenStr string) (models.AccessToken, error) {
	filter := bson.M{"accessToken": tokenStr}
	var result models.AccessToken
	collection := s.getAccessTokenCollection()
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	if result.ExpiresAt.UTC().Before(time.Now().UTC()) {
		_, err = collection.DeleteOne(context.Background(), filter)
		if err != nil {
			return result, err
		}
	}
	return result, err
}

func (s *AuthService) QueryAccessTokenByClientId(clientId string) (models.AccessToken, error) {
	filter := bson.M{"clientId": clientId}
	var result models.AccessToken
	collection := s.getAccessTokenCollection()
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// fetch an access token from the config auth endpoint
func (s *AuthService) FetchAccessTokenFromAuthEndpoint(authEndpoint string, clientId string, clientSecret string) (models.AccessToken, error) {
	tokenRequest := dto.AccessTokenRequestData{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	payload, err := json.Marshal(tokenRequest)
	if err != nil {
		return models.AccessToken{}, err
	}

	resp, err := http.Post(authEndpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return models.AccessToken{}, err
	}

	defer resp.Body.Close()

	var accessToken models.AccessToken
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.AccessToken{}, err
	}

	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return models.AccessToken{}, err
	}

	accessToken.ClientId = clientId

	return accessToken, nil
}
