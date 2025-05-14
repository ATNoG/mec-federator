package services

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
