package services

import (
	"context"
	"fmt"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

/*
 * MecSystemService
 *	responsible for interacting with a MECS
 */

type MecSystemServiceInterface interface {
	GetMecInfo() (models.OrchestratorInfo, error)
	SetMecInfo(orchestratorInfo models.OrchestratorInfo) error
}

type MecSystemService struct {
	mongoClient *mongo.Client
}

func NewMecSystemService(mongoClient *mongo.Client) *MecSystemService {
	return &MecSystemService{
		mongoClient: mongoClient,
	}
}

func (mss *MecSystemService) getMecSystemCollection() *mongo.Collection {
	return mss.mongoClient.Database(config.AppConfig.Database).Collection("systems")
}

func (mecSystemService *MecSystemService) GetMecInfo() (models.OrchestratorInfo, error) {
	collection := mecSystemService.getMecSystemCollection()

	filter := bson.M{}

	orchestratorInfo, err := utils.FetchEntityFromDatabase[models.OrchestratorInfo](collection, filter)
	if err != nil {
		return models.OrchestratorInfo{}, fmt.Errorf("error. could not fetch orchestrator info from database: %s", err)
	}

	return orchestratorInfo, nil
}

func (mecSystemService *MecSystemService) SetMecInfo(orchestratorInfo models.OrchestratorInfo) error {
	collection := mecSystemService.getMecSystemCollection()

	filter := bson.M{"operatorId": orchestratorInfo.OperatorId}

	// fields present in the payload shall overide old values
	fieldsToUpdate := bson.M{}
	if orchestratorInfo.OperatorId != "" {
		fieldsToUpdate["operatorId"] = orchestratorInfo.OperatorId
	}
	if orchestratorInfo.OperatorName != "" {
		fieldsToUpdate["operatorName"] = orchestratorInfo.OperatorName
	}
	if orchestratorInfo.OperatorCountryCode != "" {
		fieldsToUpdate["operatorCountryCode"] = orchestratorInfo.OperatorCountryCode
	}
	if orchestratorInfo.MobileCountryCode != "" {
		fieldsToUpdate["mobileCountryCode"] = orchestratorInfo.MobileCountryCode
	}
	if len(orchestratorInfo.MobileNetworkCodes) > 0 {
		fieldsToUpdate["mobileNetworkCodes"] = orchestratorInfo.MobileNetworkCodes
	}
	if orchestratorInfo.KafkaUrl != "" {
		fieldsToUpdate["kafkaUrl"] = orchestratorInfo.KafkaUrl
	}

	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": fieldsToUpdate})
	return err
}
