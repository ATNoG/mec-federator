package ewbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/services"
	"github.com/mankings/mec-federator/internal/utils"
)

type ZonesInfoSyncController struct {
	zoneService         *services.ZoneService
	orchestratorService *services.OrchestratorService
	kafkaClientService  *services.KafkaClientService
}

func NewZonesInfoSyncController(zoneService *services.ZoneService, orchestratorService *services.OrchestratorService, kafkaClientService *services.KafkaClientService) *ZonesInfoSyncController {
	return &ZonesInfoSyncController{
		zoneService:         zoneService,
		orchestratorService: orchestratorService,
		kafkaClientService:  kafkaClientService,
	}
}

// @Summary Subscribe to a Zone
// @Description Used by origin OP to show intent on using a partner OP's zone
// @Tags EWBI - ZonesInfoSync
func (zisc *ZonesInfoSyncController) SubscribeZoneController(c *gin.Context) {
	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	// get the request body and decode
	var zoneRegistrationRequestData models.ZoneRegistrationRequestData
	if err := c.ShouldBindJSON(&zoneRegistrationRequestData); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// for each string in the AcceptedAvailabilityZones array, get the zone details from the database
	acceptedZones := []models.ZoneDetails{}
	for _, zoneId := range zoneRegistrationRequestData.AcceptedAvailabilityZones {
		zone, err := zisc.zoneService.GetLocalZoneById(zoneId)
		if err != nil {
			utils.HandleProblem(c, http.StatusInternalServerError, "Error getting zone details for zoneId: "+zoneId)
			return
		}

		// check if they are already subscribed in this federation context
		_, err = zisc.zoneService.GetZoneRegisteredData(federationContextId, zoneId)
		if err != nil {
			utils.HandleProblem(c, http.StatusInternalServerError, "Error checking if zone is subscribed: "+zoneId)
			return
		}

		acceptedZones = append(acceptedZones, zone)
	}

	// make the response body
	zoneRegistrationResponseData := models.ZoneRegistrationResponseData{}

	// for each zone in the acceptedzones, make the corresponding ZoneRegisteredData
	acceptedZoneResourceInfo := []models.ZoneRegisteredData{}
	for _, zone := range acceptedZones {
		zoneRegisteredData := models.ZoneRegisteredData{
			ZoneId:              zone.ZoneId,
			FederationContextId: federationContextId,
		}

		acceptedZoneResourceInfo = append(acceptedZoneResourceInfo, zoneRegisteredData)
	}

	zoneRegistrationResponseData.AcceptedZoneResourceInfo = acceptedZoneResourceInfo

	// return the accepted zones
	c.JSON(http.StatusOK, zoneRegistrationResponseData)
}

// @Summary Unsubscribe from a Zone
// @Description Used by origin OP to show intent on not using a partner OP's zone anymore
// @Tags EWBI - ZonesInfoSync
func (zisc *ZonesInfoSyncController) UnsubscribeZoneController(c *gin.Context) {
	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	// get the zoneId from the path
	zoneId := c.Param("zoneId")

	// get the zone registered data for the given zoneId and federationContextId
	_, err := zisc.zoneService.GetZoneRegisteredData(federationContextId, zoneId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting zone registered data: "+zoneId)
		return
	}

	// delete the zone registered data for the given zoneId and federationContextId
	err = zisc.zoneService.DeleteZoneRegisteredData(federationContextId, zoneId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error deleting zone registered data: "+zoneId)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Zone unregistered successfully"})
}

// @Summary Get Zone Details
// @Description Used by origin OP to get details of a zone that belongs to a partner OP
// @Tags EWBI - ZonesInfoSync~
func (zisc *ZonesInfoSyncController) GetZoneController(c *gin.Context) {
	// get the zoneId from the path
	zoneId := c.Param("zoneId")

	// get the zone details from the database
	zone, err := zisc.zoneService.GetLocalZoneById(zoneId)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting zone details")
		return
	}

	c.JSON(http.StatusOK, zone)
}

// @Summary Get All Local Zones
// @Description Used by origin OP to get all zones that belong to a partner OP
// @Tags EWBI - ZonesInfoSync
func (zisc *ZonesInfoSyncController) GetAllLocalZonesController(c *gin.Context) {
	// ensure latest zones are up to date
	localZones, err := zisc.zoneService.GetLocalZones()
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting local zones")
		return
	}

	c.JSON(http.StatusOK, localZones)
}

func (zisc *ZonesInfoSyncController) PostMetricsController(c *gin.Context) {
	// get the request body and decode
	var metricsRequestData dto.OrchMehMetricsRequestData
	if err := c.ShouldBindJSON(&metricsRequestData); err != nil {
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// post metrics to kafka topic
	metricsTopic := "federation-meh-metrics"
	metricsMessage := metricsRequestData

	msgId, err := zisc.kafkaClientService.Produce(metricsTopic, metricsMessage)
	if err != nil {
		utils.HandleProblem(c, http.StatusInternalServerError, "Error posting metrics to kafka: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"msgId": msgId})
}
