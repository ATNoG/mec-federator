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

func (omc *MecSystemManagementController) GetOrchestratorInfoController(c *gin.Context) {
	log.Print("GetOrchestratorInfoController - Getting orchestrator information")

	orchestratorInfo, err := omc.mecSystemService.GetMecInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orchestratorInfo)
}
