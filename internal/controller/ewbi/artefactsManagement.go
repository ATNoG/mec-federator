package ewbi

import (
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type ArtefactManagementController struct {
	orchestratorService *services.OrchestratorService
	artefactService     *services.ArtefactService
}

func NewArtefactManagementController(orchestratorService *services.OrchestratorService, artefactService *services.ArtefactService) *ArtefactManagementController {
	return &ArtefactManagementController{
		orchestratorService: orchestratorService,
		artefactService:     artefactService,
	}
}

// @Summary Onboard an artefact
// @Description Receives an artefact from origin OP. Artefact is a zip file containing scripts and/or packaging files
// @Tags EWBI - ArtefactManagement
func (amc *ArtefactManagementController) OnboardArtefactController(c *gin.Context) {
	log.Print("OnboardArtefactController - Onboarding artefact onto federator")

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Could not parse multipart form")
		return
	}

	// Bind the form to the artefactOnboardRequest
	var artefactOnboardRequest dto.ArtefactOnboardRequest
	if err := c.ShouldBind(&artefactOnboardRequest); err != nil {
		slog.Error("Could not bind form to artefactOnboardRequest", "error", err)
		utils.HandleProblem(c, http.StatusBadRequest, "Could not bind form to artefactOnboardRequest")
		return
	}

	// Validate the artefactOnboardRequest
	if err := artefactOnboardRequest.Validate(); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Artefact onboarding request", "artefactOnboardRequest", artefactOnboardRequest)

	// Get the file from the form
	artefactFiles := form.File["artefactFile"]
	if len(artefactFiles) == 0 {
		utils.HandleProblem(c, http.StatusBadRequest, "Missing artefact file")
		return
	}

	// Read the file content
	file, err := artefactFiles[0].Open()
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error reading artefact file")
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error reading artefact file content")
		return
	}

	// // Validate that file is a valid descriptor file
	// if !utils.IsDescriptorFile(fileContent) {
	// 	utils.HandleProblem(c, http.StatusBadRequest, "Invalid descriptor file")
	// 	return
	// }

	// Create artefact object
	artefact := models.Artefact{
		Id:                  artefactOnboardRequest.ArtefactId,
		FederationContextId: c.GetString("federationContextId"),
		Name:                artefactOnboardRequest.ArtefactName,
		Description:         artefactOnboardRequest.ArtefactDescription,
		VirtType:            artefactOnboardRequest.ArtefactVirtType,
		DescriptorType:      artefactOnboardRequest.ArtefactDescriptorType,
		FileName:            artefactOnboardRequest.ArtefactFileName,
		FileFormat:          artefactOnboardRequest.ArtefactFileFormat,
		ArtefactFile:        &fileContent,
	}

	// Save artefact object to database
	err = amc.artefactService.SaveArtefact(artefact)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error saving artefact to database")
		return
	}

	// Onboard the artefact onto the orchestrator
	err = amc.orchestratorService.OnboardArtefact(artefact)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error onboarding artefact onto orchestrator")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "Artefact onboarded successfully",
		"artefactId": artefactOnboardRequest.ArtefactId,
	})
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
