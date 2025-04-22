package services

import (
	"context"
	"log"
	"time"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// interface for authentication service
type AuthServiceInterface interface {
	SaveAccessToken(accessToken models.AccessToken) error
	QueryAccessToken(tokenStr string) (models.AccessToken, error)
}

// implementation of authentication service, requiring a mongo client
type AuthService struct {
	mongoClient *mongo.Client
}

// create a new authentication service, injecting given mongo client
func NewAuthService(mongoClient *mongo.Client) *AuthService {
	return &AuthService{mongoClient: mongoClient}
}

// save an access token to the database
func (s *AuthService) SaveAccessToken(accessToken models.AccessToken) error {
	collection := s.mongoClient.Database(config.AppConfig.Database).Collection("access_tokens")
	_, err := collection.InsertOne(context.TODO(), accessToken)
	log.Printf("AuthService - Saved access token")
	return err
}

// query an access token from the database; if experired, delete it
func (s *AuthService) QueryAccessToken(tokenStr string) (models.AccessToken, error) {
	filter := bson.M{"accessToken": tokenStr}
	var result models.AccessToken
	collection := s.mongoClient.Database(config.AppConfig.Database).Collection("access_tokens")
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
