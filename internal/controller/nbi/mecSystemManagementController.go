package nbi

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/services"
)

type MecSystemManagementController struct {
	orchestratorService *services.OrchestratorService
	mecSystemService    *services.MecSystemService
}

func NewMecSystemManagementController(orchestratorService *services.OrchestratorService, mecSystemService *services.MecSystemService) *MecSystemManagementController {
	return &MecSystemManagementController{
		orchestratorService: orchestratorService,
		mecSystemService:    mecSystemService,
	}
}

// @Summary Register MEC System Information
// @Description Provide the MEF with info about the MEC System where it belongs to
// @Tags NBI - MEC System Management
// @Accept json
// @Produce json
// @Param orchestratorInfo body models.OrchestratorInfo true "Orchestrator Information"
// @Success 200 {object} models.OrchestratorInfo "Orchestrator Information"
// @Failure 400 {object} models.ProblemDetails "Invalid request body or missing fields"
// @Failure 500 {object} models.ProblemDetails "Internal error during orchestrator registration"
// @Router /nbi/mec-system [post]
func (omc *MecSystemManagementController) RegisterOrchestratorController(c *gin.Context) {
	log.Print("RegisterOrchestratorController - Registering orchestrator information")

	var orchestratorInfo models.OrchestratorInfo
	if err := c.ShouldBindJSON(&orchestratorInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := omc.mecSystemService.SetMecInfo(orchestratorInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Orchestrator information registered successfully"})
}

// @Summary Get MEC System Information
// @Description Get the MEC System Information from the MEF
// @Tags NBI - MEC System Management
// @Accept json
// @Produce json
// @Success 200 {object} models.OrchestratorInfo "Orchestrator Information"
// @Failure 500 {object} models.ProblemDetails "Internal error during orchestrator information retrieval"
// @Router /nbi/mec-system [get]
func (omc *MecSystemManagementController) GetOrchestratorInfoController(c *gin.Context) {
	log.Print("GetOrchestratorInfoController - Getting orchestrator information")

	orchestratorInfo, err := omc.mecSystemService.GetMecInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orchestratorInfo)
}
