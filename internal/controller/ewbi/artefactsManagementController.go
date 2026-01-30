package ewbi

import (
	"io"
	"log"
	"net/http"

	"github.com/atnog/mec-federator/internal/models"
	"github.com/atnog/mec-federator/internal/models/dto"
	"github.com/atnog/mec-federator/internal/services"
	"github.com/atnog/mec-federator/internal/utils"
	"github.com/gin-gonic/gin"
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

// @Summary Onboard Artefact
// @Description Receives and onboards an artefact from the origin operator. The artefact is a packaged file (zip/tar) containing application scripts and packaging files. The artefact is validated and stored in the database.
// @Tags EWBI - ArtefactManagement
// @Accept multipart/form-data
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param artefactFile formData file true "Artefact package file (zip/tar format)"
// @Param artefactId formData string true "Unique identifier for the artefact"
// @Param appProviderId formData string true "Application provider identifier"
// @Param artefactName formData string true "Human-readable name of the artefact"
// @Param artefactVersionInfo formData string true "Version information of the artefact"
// @Param artefactDescription formData string true "Detailed description of the artefact"
// @Param artefactVirtType formData string true "Virtualization type (e.g., container, VM)"
// @Param artefactDescriptorType formData string true "Descriptor file type (e.g., TOSCA, Helm)"
// @Param artefactFileFormat formData string true "File format of the artefact (e.g., zip, tar.gz)"
// @Param artefactFileName formData string true "Original filename of the artefact"
// @Success 200 {object} map[string]string "status: Artefact onboarded successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid form data or validation errors"
// @Failure 500 {object} models.ProblemDetails "File processing or database errors"
// @Router /ewbi/{federationContextId}/artefact [post]
func (amc *ArtefactManagementController) OnboardArtefactController(c *gin.Context) {
	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-onboard-artefact-init",
		Message: "",
		Value:   nil,
	})

	federationContextId := c.Param("federationContextId")
	log.Printf("OnboardArtefactController - Starting artefact onboarding for federation: %s", federationContextId)

	// Parse the multipart form
	log.Printf("OnboardArtefactController - Parsing multipart form for federation: %s", federationContextId)
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("OnboardArtefactController - Error parsing multipart form for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusBadRequest, "Error parsing multipart form: "+err.Error())
		return
	}

	// Bind the form to the artefactOnboardRequest
	log.Printf("OnboardArtefactController - Binding form values for federation: %s", federationContextId)
	var artefactOnboardRequest dto.ArtefactOnboardRequest
	if err := c.ShouldBind(&artefactOnboardRequest); err != nil {
		log.Printf("OnboardArtefactController - Error binding form values for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusBadRequest, "Error binding form values: "+err.Error())
		return
	}

	// Validate the artefactOnboardRequest
	log.Printf("OnboardArtefactController - Validating artefact onboard request for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
	if err := artefactOnboardRequest.Validate(); err != nil {
		log.Printf("OnboardArtefactController - Validation failed for artefact %s in federation %s: %v", artefactOnboardRequest.ArtefactId, federationContextId, err)
		utils.HandleProblem(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the file from the form
	log.Printf("OnboardArtefactController - Extracting artefact file for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
	artefactFiles := form.File["artefactFile"]
	if len(artefactFiles) == 0 {
		log.Printf("OnboardArtefactController - Missing artefact file for federation %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
		utils.HandleProblem(c, http.StatusBadRequest, "Missing artefact file")
		return
	}

	// Read the file content
	log.Printf("OnboardArtefactController - Reading artefact file content for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
	file, err := artefactFiles[0].Open()
	if err != nil {
		log.Printf("OnboardArtefactController - Error opening artefact file for federation %s, artefactId %s: %v", federationContextId, artefactOnboardRequest.ArtefactId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error reading file: "+err.Error())
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Printf("OnboardArtefactController - Error reading file content for federation %s, artefactId %s: %v", federationContextId, artefactOnboardRequest.ArtefactId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error reading file content: "+err.Error())
		return
	}

	// // Get descriptor data from the tar ball
	// log.Printf("OnboardArtefactController - Extracting descriptor data for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
	// descriptorData, err := utils.GetDescriptorData(fileContent)
	// if err != nil {
	// 	log.Printf("OnboardArtefactController - Error getting descriptor data for federation %s, artefactId %s: %v", federationContextId, artefactOnboardRequest.ArtefactId, err)
	// 	utils.HandleProblem(c, http.StatusBadRequest, "Error getting descriptor data: "+err.Error())
	// 	return
	// }

	// // Validate the descriptor data
	// log.Printf("OnboardArtefactController - Validating descriptor data for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
	// appPkg, err := utils.ValidateDescriptorData(descriptorData)
	// if err != nil {
	// 	log.Printf("OnboardArtefactController - Error validating descriptor data for federation %s, artefactId %s: %v", federationContextId, artefactOnboardRequest.ArtefactId, err)
	// 	utils.HandleProblem(c, http.StatusBadRequest, "Error validating descriptor data: "+err.Error())
	// 	return
	// }
	// appPkg.AppD = fileContent

	// // Onboard the artefact onto the orchestrator
	// log.Printf("OnboardArtefactController - Onboarding artefact to orchestrator for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
	// appPkgId, err := amc.orchestratorService.OnboardAppPkg(appPkg)
	// if err != nil {
	// 	log.Printf("OnboardArtefactController - Error onboarding artefact to orchestrator for federation %s, artefactId %s: %v", federationContextId, artefactOnboardRequest.ArtefactId, err)
	// 	utils.HandleProblem(c, http.StatusInternalServerError, "Error onboarding artefact onto orchestrator: "+err.Error())
	// 	return
	// }

	// log.Printf("OnboardArtefactController - Artefact onboarded to orchestrator successfully for federation: %s, artefactId: %s, appPkgId: %s", federationContextId, artefactOnboardRequest.ArtefactId, appPkgId)

	// Create artefact object
	log.Printf("OnboardArtefactController - Creating artefact object for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
	artefact := models.Artefact{
		Id:                  artefactOnboardRequest.ArtefactId,
		FederationContextId: federationContextId,
		AppProviderId:       artefactOnboardRequest.AppProviderId,
		Name:                artefactOnboardRequest.ArtefactName,
		VersionInfo:         artefactOnboardRequest.ArtefactVersionInfo,
		Description:         artefactOnboardRequest.ArtefactDescription,
		VirtType:            artefactOnboardRequest.ArtefactVirtType,
		DescriptorType:      artefactOnboardRequest.ArtefactDescriptorType,
		FileName:            artefactOnboardRequest.ArtefactFileName,
		FileFormat:          artefactOnboardRequest.ArtefactFileFormat,
		ArtefactFile:        &fileContent,
	}

	// Save artefact object to database
	log.Printf("OnboardArtefactController - Saving artefact to database for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)
	err = amc.artefactService.SaveArtefact(artefact)
	if err != nil {
		log.Printf("OnboardArtefactController - Error saving artefact to database for federation %s, artefactId %s: %v", federationContextId, artefactOnboardRequest.ArtefactId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error saving artefact to database: "+err.Error())
		return
	}

	log.Printf("OnboardArtefactController - Artefact onboarded successfully for federation: %s, artefactId: %s", federationContextId, artefactOnboardRequest.ArtefactId)

	utils.SendResultsMessage(utils.ResultsMessage{
		Name:    "federation-po-onboard-artefact-done",
		Message: "",
		Value:   nil,
	})

	c.JSON(http.StatusOK, gin.H{"status": "Artefact onboarded successfully"})
}

// @Summary Get Artefact Details
// @Description Retrieves detailed information about a specific artefact by its ID within a federation context
// @Tags EWBI - ArtefactManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param artefactId path string true "Artefact identifier"
// @Success 200 {object} dto.GetArtefactResponse
// @Failure 404 {object} models.ProblemDetails "Artefact not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /ewbi/{federationContextId}/artefacts/{artefactId} [get]
func (amc *ArtefactManagementController) GetArtefactController(c *gin.Context) {
	// get the artefact id from path
	artefactId := c.Param("artefactId")

	// get the federation context id from path
	federationContextId := c.Param("federationContextId")

	log.Printf("GetArtefactController - Getting artefact details for federation: %s, artefactId: %s", federationContextId, artefactId)

	// get the artefact from the database
	log.Printf("GetArtefactController - Retrieving artefact from database for federation: %s, artefactId: %s", federationContextId, artefactId)
	artefact, err := amc.artefactService.GetArtefact(federationContextId, artefactId)
	if err != nil {
		log.Printf("GetArtefactController - Artefact not found for federation %s, artefactId %s: %v", federationContextId, artefactId, err)
		utils.HandleProblem(c, http.StatusNotFound, "Artefact not found: "+err.Error())
		return
	}

	// build the response
	log.Printf("GetArtefactController - Building response for federation: %s, artefactId: %s", federationContextId, artefactId)
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

	log.Printf("GetArtefactController - Successfully retrieved artefact details for federation: %s, artefactId: %s", federationContextId, artefactId)
	c.JSON(http.StatusOK, response)
}

// @Summary Delete Artefact
// @Description Removes an artefact from the database. Artefacts that have been onboarded to the orchestrator cannot be deleted until the onboarding is removed.
// @Tags EWBI - ArtefactManagement
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param artefactId path string true "Artefact identifier"
// @Success 200 {object} map[string]string "status: Artefact deleted successfully"
// @Failure 403 {object} models.ProblemDetails "Artefact has been onboarded and cannot be deleted"
// @Failure 404 {object} models.ProblemDetails "Artefact not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /ewbi/{federationContextId}/artefact/{artefactId} [delete]
func (amc *ArtefactManagementController) DeleteArtefactController(c *gin.Context) {
	// get the artefact id from path
	artefactId := c.Param("artefactId")

	// get the federation context id from path
	federationContextId := c.Param("federationContextId")

	log.Printf("DeleteArtefactController - Starting artefact deletion for federation: %s, artefactId: %s", federationContextId, artefactId)

	// check if the artefact exists
	log.Printf("DeleteArtefactController - Checking if artefact exists for federation: %s, artefactId: %s", federationContextId, artefactId)
	artefact, err := amc.artefactService.GetArtefact(federationContextId, artefactId)
	if err != nil {
		log.Printf("DeleteArtefactController - Artefact not found for federation %s, artefactId %s: %v", federationContextId, artefactId, err)
		utils.HandleProblem(c, http.StatusNotFound, "Artefact not found: "+err.Error())
		return
	}

	// check if appPkgId has a value
	if artefact.AppPkgId != "" {
		log.Printf("DeleteArtefactController - Cannot delete artefact with AppPkgId for federation %s, artefactId %s", federationContextId, artefactId)
		utils.HandleProblem(c, http.StatusForbidden, "Cannot delete artefact that has been onboarded to orchestrator. Delete onboarding first.")
		return
	}

	// delete artefact from the database
	log.Printf("DeleteArtefactController - Deleting artefact from database for federation: %s, artefactId: %s", federationContextId, artefactId)
	err = amc.artefactService.RemoveArtefact(federationContextId, artefactId)
	if err != nil {
		log.Printf("DeleteArtefactController - Error deleting artefact from database for federation %s, artefactId %s: %v", federationContextId, artefactId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error deleting artefact from database: "+err.Error())
		return
	}

	log.Printf("DeleteArtefactController - Artefact deleted successfully for federation: %s, artefactId: %s", federationContextId, artefactId)
	c.JSON(http.StatusOK, gin.H{"status": "Artefact deleted successfully"})
}

func (amc *ArtefactManagementController) UploadFileController(c *gin.Context) {

}

func (amc *ArtefactManagementController) GetFileController(c *gin.Context) {

}

func (amc *ArtefactManagementController) DeleteFileController(c *gin.Context) {

}
