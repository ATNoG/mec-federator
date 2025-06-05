package services

import (
	"context"
	"time"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AppInstanceServiceInterface interface {
	RegisterAppInstance(appInstance models.AppInstance) error
}

type AppInstanceService struct {
	KafkaService *KafkaService
}

func NewAppInstanceService(kafkaService *KafkaService) *AppInstanceService {
	return &AppInstanceService{
		KafkaService: kafkaService,
	}
}

func (ais *AppInstanceService) getAppInstanceCollection() *mongo.Collection {
	return config.GetMongoDatabase().Collection("appInstances")
}

// Registers an app instance to the database
func (ais *AppInstanceService) RegisterAppInstance(appInstance models.AppInstance) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := ais.getAppInstanceCollection()
	_, err := collection.InsertOne(ctx, appInstance)
	return err
}

// Provides a list of instances for a given federationContextId
func (ais *AppInstanceService) GetAppInstancesByFederationContextId(fedContextId string) error {
	return nil
}
