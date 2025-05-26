package services

import (
	"context"
	"log/slog"

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
	_, err = s.kafkaService.Produce("new_app_pkg", message)
	if err != nil {
		return "", err
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
	_, err := s.kafkaService.Produce("delete_app_pkg", message)
	if err != nil {
		return err
	}

	return nil
}
