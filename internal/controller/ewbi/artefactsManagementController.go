package ewbi

import (
	"io"
	"log"
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
// @Param federationContextId path string true "Federation Context ID"
// @Accept multipart/form-data
// @Produce json
// @Param artefactFile formData file true "Artefact file"
// @Param artefactId formData string true "Artefact ID"
// @Param appProviderId formData string true "App Provider ID"
// @Param artefactName formData string true "Artefact Name"
// @Param artefactVersionInfo formData string true "Artefact Version Info"
// @Param artefactDescription formData string true "Artefact Description"
// @Param artefactVirtType formData string true "Artefact Virt Type"
// @Param artefactDescriptorType formData string true "Artefact Descriptor Type"
// @Param artefactFileFormat formData string true "Artefact File Format"
// @Param artefactFileName formData string true "Artefact File Name"
// @Success 200 {object} map[string]string "status: Artefact onboarded successfully"
// @Failure 400 {object} models.ProblemDetails
// @Failure 500 {object} models.ProblemDetails
// @Router /ewbi/{federationContextId}/artefact [post]
func (amc *ArtefactManagementController) OnboardArtefactController(c *gin.Context) {
	log.Print("OnboardArtefactController - Onboarding artefact onto federator")

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Error parsing multipart form: "+err.Error())
		return
	}

	// Bind the form to the artefactOnboardRequest
	var artefactOnboardRequest dto.ArtefactOnboardRequest
	if err := c.ShouldBind(&artefactOnboardRequest); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Error binding form values: "+err.Error())
		return
	}

	// Validate the artefactOnboardRequest
	if err := artefactOnboardRequest.Validate(); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the file from the form
	artefactFiles := form.File["artefactFile"]
	if len(artefactFiles) == 0 {
		utils.HandleProblem(c, http.StatusBadRequest, "Missing artefact file")
		return
	}

	// Read the file content
	file, err := artefactFiles[0].Open()
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error reading file: "+err.Error())
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error reading file content: "+err.Error())
		return
	}

	// Get descriptor data from the tar ball
	descriptorData, err := utils.GetDescriptorData(fileContent)
	if err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Error getting descriptor data: "+err.Error())
		return
	}

	// Validate the descriptor data
	appPkg, err := utils.ValidateDescriptorData(descriptorData)
	if err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Error validating descriptor data: "+err.Error())
		return
	}
	appPkg.AppD = fileContent

	// Onboard the artefact onto the orchestrator
	appPkgId, err := amc.orchestratorService.OnboardAppPkg(appPkg)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error onboarding artefact onto orchestrator: "+err.Error())
		return
	}

	// Create artefact object
	artefact := models.Artefact{
		Id:                  artefactOnboardRequest.ArtefactId,
		FederationContextId: c.GetString("federationContextId"),
		AppProviderId:       artefactOnboardRequest.AppProviderId,
		Name:                artefactOnboardRequest.ArtefactName,
		VersionInfo:         artefactOnboardRequest.ArtefactVersionInfo,
		Description:         artefactOnboardRequest.ArtefactDescription,
		VirtType:            artefactOnboardRequest.ArtefactVirtType,
		DescriptorType:      artefactOnboardRequest.ArtefactDescriptorType,
		FileName:            artefactOnboardRequest.ArtefactFileName,
		FileFormat:          artefactOnboardRequest.ArtefactFileFormat,
		ArtefactFile:        &fileContent,
		AppPkgId:            appPkgId,
	}

	// Save artefact object to database
	err = amc.artefactService.SaveArtefact(artefact)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error saving artefact to database: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Artefact onboarded successfully"})
}

// @Summary Get an artefact
// @Description Get an artefact details by its id
// @Tags EWBI - ArtefactManagement
// @Param federationContextId path string true "Federation Context ID"
// @Param artefactId path string true "Artefact ID"
// @Success 200 {object} dto.GetArtefactResponse
// @Failure 404 {object} models.ProblemDetails
// @Failure 500 {object} models.ProblemDetails
// @Router /ewbi/{federationContextId}/artefacts/{artefactId} [get]
func (amc *ArtefactManagementController) GetArtefactController(c *gin.Context) {
	log.Print("GetArtefactController - Getting artefact details")

	// get the artefact id from path
	artefactId := c.Param("artefactId")

	// get the federation context id from path
	federationContextId := c.Param("federationContextId")

	// get the artefact from the database
	artefact, err := amc.artefactService.GetArtefact(federationContextId, artefactId)
	if err != nil {
		utils.HandleProblem(c, http.StatusNotFound, "Artefact not found: "+err.Error())
		return
	}

	// build the response
	response := dto.GetArtefactResponse{
		ArtefactId:             artefact.Id,
		AppProviderId:          artefact.AppProviderId,
		ArtefactVersionInfo:    artefact.VersionInfo,
		ArtefactName:           artefact.Name,
		ArtefactDescription:    artefact.Description,
		ArtefactVirtType:       artefact.VirtType,
		ArtefactDescriptorType: artefact.DescriptorType,
		ArtefactFileFormat:     artefact.FileFormat,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Delete an artefact
// @Description Delete an artefact by its id
// @Tags EWBI - ArtefactManagement
// @Param federationContextId path string true "Federation Context ID"
// @Param artefactId path string true "Artefact ID"
// @Success 200 {object} map[string]string "status: Artefact deleted successfully"
// @Failure 404 {object} models.ProblemDetails
// @Failure 500 {object} models.ProblemDetails
// @Router /ewbi/{federationContextId}/artefact/{artefactId} [delete]
func (amc *ArtefactManagementController) DeleteArtefactController(c *gin.Context) {
	log.Print("DeleteArtefactController - Deleting artefact")

	// get the artefact id from path
	artefactId := c.Param("artefactId")

	// get the federation context id from path
	federationContextId := c.Param("federationContextId")

	// check if the artefact exists
	artefact, err := amc.artefactService.GetArtefact(federationContextId, artefactId)
	if err != nil {
		utils.HandleProblem(c, http.StatusNotFound, "Artefact not found: "+err.Error())
		return
	}

	// remove artefact from the orchestrator
	err = amc.orchestratorService.RemoveAppPkg(artefact.AppPkgId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error removing artefact from orchestrator: "+err.Error())
		return
	}

	// delete artefact from the database
	err = amc.artefactService.RemoveArtefact(federationContextId, artefactId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error deleting artefact from database: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Artefact deleted successfully"})
}

func (amc *ArtefactManagementController) UploadFileController(c *gin.Context) {

}

func (amc *ArtefactManagementController) GetFileController(c *gin.Context) {

}

func (amc *ArtefactManagementController) DeleteFileController(c *gin.Context) {

}
