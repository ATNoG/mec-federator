package services

import (
	"context"
	"time"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ApplicationServiceInterface interface {
	CreateApplication(application models.Application) error
	GetApplication(federationContextId string, appId string) (models.Application, error)
	UpdateApplication(application models.Application) error
	DeleteApplication(federationContextId string, appId string) error
}

type ApplicationService struct {
	kafkaService *KafkaClientService
}

func NewApplicationService(kafkaService *KafkaClientService) *ApplicationService {
	return &ApplicationService{
		kafkaService: kafkaService,
	}
}

func (as *ApplicationService) getApplicationCollection() *mongo.Collection {
	return config.GetMongoDatabase().Collection("applications")
}

func (as *ApplicationService) CreateApplication(application models.Application) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := as.getApplicationCollection()
	_, err := collection.InsertOne(ctx, application)
	return err
}

func (as *ApplicationService) GetApplication(federationContextId string, appId string) (models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := as.getApplicationCollection()
	filter := bson.M{"federationContextId": federationContextId, "appId": appId}
	var application models.Application
	err := collection.FindOne(ctx, filter).Decode(&application)
	return application, err
}

func (as *ApplicationService) DeleteApplication(federationContextId string, appId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := as.getApplicationCollection()
	filter := bson.M{"federationContextId": federationContextId, "appId": appId}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}
