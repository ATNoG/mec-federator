package ewbi

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/services"
)

type ArtefactManagementController struct {
	orchestratorService *services.OrchestratorService
}

func NewArtefactManagementController(orchestratorService *services.OrchestratorService) *ArtefactManagementController {
	return &ArtefactManagementController{
		orchestratorService: orchestratorService,
	}
}

// @Summary Onboard an artefact
// @Description Receives an artefact from origin OP. Artefact is a zip file containing scripts and/or packaging files
// @Tags EWBI - ArtefactManagement
func (amc *ArtefactManagementController) OnboardArtefactController(c *gin.Context) {
	log.Print("OnboardArtefactController - Onboarding artefact onto federator")

	// Check if the request multipart form is valid
	
}

func (amc *ArtefactManagementController) GetArtefactController(c *gin.Context) {

}

func (amc *ArtefactManagementController) DeleteArtefactController(c *gin.Context) {

}

func (amc *ArtefactManagementController) UploadFileController(c *gin.Context) {

}

func (amc *ArtefactManagementController) GetFileController(c *gin.Context) {

}

func (amc *ArtefactManagementController) DeleteFileController(c *gin.Context) {

}
