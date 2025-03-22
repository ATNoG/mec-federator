package services

/*
 * OrchestratorService
 *	responsible for interacting with an orchestrator
 */

type OrchestratorService interface {
}

type OrchestratorServiceImpl struct {
}

func NewOrchestratorService() *OrchestratorServiceImpl {
	return &OrchestratorServiceImpl{}
}
