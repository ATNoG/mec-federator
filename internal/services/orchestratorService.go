package services

/*
 * OrchestratorService
 *	responsible for interacting with the registered orchestrator
 */

type OrchestratorServiceInterface interface {
}

type OrchestratorService struct {
}

func NewOrchestratorService() *OrchestratorService {
	return &OrchestratorService{}
}
