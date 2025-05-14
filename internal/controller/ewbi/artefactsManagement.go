package ewbi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
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
	form, err := c.MultipartForm()
	if err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body, missing parameters or wrong data type")
		return
	}
	artefactFile := form.File["artefactFile"]
	if len(artefactFile) == 0 {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body, missing parameters or wrong data type")
		return
	}
	

	// Check if the artefact is valid
	// Check if the artefact is already onboarded
	// Save the artefact to the database
	// Send message to orchestrator through kafka with artefact id

}

func (amc *ArtefactManagementController) GetArtefactController(c *gin.Context) {
	log.Print("GetArtefactController - Getting artefact details")
}

func (amc *ArtefactManagementController) DeleteArtefactController(c *gin.Context) {

}

func (amc *ArtefactManagementController) UploadFileController(c *gin.Context) {

}

func (amc *ArtefactManagementController) GetFileController(c *gin.Context) {

}

func (amc *ArtefactManagementController) DeleteFileController(c *gin.Context) {

}
