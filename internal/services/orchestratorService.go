package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

/*
 * OrchestratorService
 *	responsible for interacting with the registered orchestrator
 */

type OrchestratorServiceInterface interface {
}

type OrchestratorService struct {
	kafkaService *KafkaService
}

func NewOrchestratorService(kafkaService *KafkaService) *OrchestratorService {
	return &OrchestratorService{
		kafkaService: kafkaService,
	}
}

func (s *OrchestratorService) getOrchestratorAppPkgsCollection() *mongo.Collection {
	return config.GetOrchestratorMongoDatabase().Collection("app_pkgs")
}

// Onboards an artefact onto the orchestrator
func (s *OrchestratorService) OnboardAppPkg(appPkg dto.NewAppPkg) (string, error) {
	slog.Info("Onboarding appPkg onto orchestrator")

	// Insert the app_pkg into the database
	result, err := s.getOrchestratorAppPkgsCollection().InsertOne(context.Background(), appPkg)
	if err != nil {
		return "", err
	}

	// get the app_pkg id as a string
	appPkgId := result.InsertedID.(bson.ObjectID)

	// make a message to send to the kafka topic
	message := dto.NewAppPkgMessage{
		AppPkgId: appPkgId.Hex(),
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaService.Produce("new_app_pkg", message)
	if err != nil {
		return "", err
	}

	// wait for a response
	rsp, err := s.kafkaService.WaitForResponse(msgId, 10*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return "", err
	}

	// check status of the response
	status := rsp["status"].(float64)
	if status != 201 {
		return "", errors.New("failed to onboard appPkg")
	}

	return appPkgId.Hex(), nil
}

// Remove an appPkg from the orchestrator
func (s *OrchestratorService) RemoveAppPkg(appPkgId string) error {
	slog.Info("Removing appPkg from orchestrator", "appPkgId", appPkgId)

	// remove the appPkg from the database
	message := dto.DeleteAppPkgMessage{
		AppPkgId: appPkgId,
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaService.Produce("delete_app_pkg", message)
	if err != nil {
		return err
	}

	// wait for a response
	rsp, err := s.kafkaService.WaitForResponse(msgId, 10*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return nil
	}

	// get status field from response
	status := rsp["status"].(float64)

	// if status is not 204, return an error
	if status != 204 {
		return errors.New("failed to remove appPkg from orchestrator")
	}

	// convert the string id to a bson.ObjectID
	appPkgIdBson, err := bson.ObjectIDFromHex(appPkgId)
	if err != nil {
		return err
	}

	// delete the appPkg from the database
	result, err := s.getOrchestratorAppPkgsCollection().DeleteOne(context.Background(), bson.M{"_id": appPkgIdBson})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("appPkg not found")
	}

	slog.Info("appPkg deleted from database", "result", result)
	return nil
}

// Instantiate an appPkg
func (s *OrchestratorService) InstantiateAppPkg(appPkgId string) error {
	slog.Info("Instantiating appPkg", "appPkgId", appPkgId)

	// make a message to send to the kafka topic
	message := dto.NewAppInstanceMessage{
		AppPkgId:    appPkgId,
		Name:        "test",
		Description: "test",
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaService.Produce("instantiate_app_pkg", message)
	if err != nil {
		return err
	}

	// wait for a response
	rsp, err := s.kafkaService.WaitForResponse(msgId, 10*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return err
	}

	// get status field from response
	status := rsp["status"].(float64)

	// if status is not 201, return an error
	if status != 201 {
		return errors.New("failed to instantiate appPkg")
	}

	return nil
}
