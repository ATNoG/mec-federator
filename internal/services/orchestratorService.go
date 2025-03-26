package services

/*
 * OrchestratorService
 *	responsible for interacting with an orchestrator
 */

type OrchestratorServiceInterface interface {
}

type OrchestratorService struct {
}

func NewOrchestratorService() *OrchestratorService {
	return &OrchestratorService{}
}

// DeleteArtefact deletes an artefact from the orchestrator
func (os *OrchestratorService) DeleteArtefact() error {
	return nil
}
