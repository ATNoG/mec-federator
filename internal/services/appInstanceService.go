package services

import (
	"context"
	"time"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

/*
 * AppInstanceService
 *	responsible for managing app instances
 */

type AppInstanceServiceInterface interface {
	RegisterAppInstance(appInstance models.AppInstance) error
}

type AppInstanceService struct {
	kafkaClientService *KafkaClientService
}

func NewAppInstanceService(kafkaClientService *KafkaClientService) *AppInstanceService {
	return &AppInstanceService{
		kafkaClientService: kafkaClientService,
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

// Removes an app instance from the database
func (ais *AppInstanceService) RemoveAppInstance(federationContextId string, appInstanceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := ais.getAppInstanceCollection()
	_, err := collection.DeleteOne(ctx, bson.M{"federationContextId": federationContextId, "appInstanceId": appInstanceId})
	return err
}

// Gets an app instance from the database
func (ais *AppInstanceService) GetAppInstance(federationContextId string, appInstanceId string) (models.AppInstance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := ais.getAppInstanceCollection()
	filter := bson.M{"federationContextId": federationContextId, "appInstanceId": appInstanceId}
	var appInstance models.AppInstance
	err := collection.FindOne(ctx, filter).Decode(&appInstance)
	if err != nil {
		return models.AppInstance{}, err
	}
	return appInstance, nil
}

// Provides a list of instances for a given federationContextId
func (ais *AppInstanceService) GetAppInstancesByFederationContextId(fedContextId string) ([]models.AppInstance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := ais.getAppInstanceCollection()
	filter := bson.M{"federationContextId": fedContextId}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var appInstances []models.AppInstance
	for cursor.Next(ctx) {
		var appInstance models.AppInstance
		err = cursor.Decode(&appInstance)
		if err != nil {
			return nil, err
		}
		appInstances = append(appInstances, appInstance)
	}
	return appInstances, nil
}

// Provides a list of all app instances resulting from any federation
func (ais *AppInstanceService) GetAllAppInstances() ([]models.AppInstance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := ais.getAppInstanceCollection()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var appInstances []models.AppInstance
	for cursor.Next(ctx) {
		var appInstance models.AppInstance
		err = cursor.Decode(&appInstance)
		if err != nil {
			return nil, err
		}
		appInstances = append(appInstances, appInstance)
	}
	return appInstances, nil
}

