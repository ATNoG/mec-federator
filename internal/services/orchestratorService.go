package services

import (
	"log/slog"

	"github.com/mankings/mec-federator/internal/models"
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

// Onboards an artefact onto the orchestrator
func (s *OrchestratorService) OnboardArtefact(artefact models.Artefact) error {
	slog.Info("Onboarding artefact onto orchestrator", "artefact", artefact)
	
	
	
	return nil
}
