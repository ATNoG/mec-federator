package ewbi

import (
	"log"
	"net/http"

	"github.com/atnog/mec-federator/internal/models"
	"github.com/atnog/mec-federator/internal/models/dto"
	"github.com/atnog/mec-federator/internal/services"
	"github.com/atnog/mec-federator/internal/utils"
	"github.com/gin-gonic/gin"
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

// @Summary Subscribe to Zone
// @Description Registers the origin operator's intent to use one or more availability zones from the partner operator
// @Tags EWBI - ZonesInfoSync
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param zoneRegistrationRequestData body models.ZoneRegistrationRequestData true "Zone Registration Request Data"
// @Success 200 {object} models.ZoneRegistrationResponseData
// @Failure 400 {object} models.ProblemDetails "Invalid request body"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /ewbi/{federationContextId}/zones [post]
func (zisc *ZonesInfoSyncController) SubscribeZoneController(c *gin.Context) {
	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")
	log.Printf("SubscribeZoneController - Starting zone subscription for federation: %s", federationContextId)

	// get the request body and decode
	log.Printf("SubscribeZoneController - Binding request body for federation: %s", federationContextId)
	var zoneRegistrationRequestData models.ZoneRegistrationRequestData
	if err := c.ShouldBindJSON(&zoneRegistrationRequestData); err != nil {
		log.Printf("SubscribeZoneController - Error binding request body for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("SubscribeZoneController - Request body: %+v", zoneRegistrationRequestData)

	// for each string in the AcceptedAvailabilityZones array, get the zone details from the database
	log.Printf("SubscribeZoneController - Processing %d zone(s) for federation: %s", len(zoneRegistrationRequestData.AcceptedAvailabilityZones), federationContextId)
	acceptedZones := []models.ZoneDetails{}
	for _, zoneId := range zoneRegistrationRequestData.AcceptedAvailabilityZones {
		log.Printf("SubscribeZoneController - Getting zone details for federation: %s, zoneId: %s", federationContextId, zoneId)
		zone, err := zisc.zoneService.GetLocalZoneById(zoneId)
		if err != nil {
			log.Printf("SubscribeZoneController - Error getting zone details for federation %s, zoneId %s: %v", federationContextId, zoneId, err)
			utils.HandleProblem(c, http.StatusInternalServerError, "Error getting zone details for zoneId: "+zoneId)
			return
		}

		// check if they are already subscribed in this federation context
		log.Printf("SubscribeZoneController - Checking if zone is already subscribed for federation: %s, zoneId: %s", federationContextId, zoneId)
		_, err = zisc.zoneService.GetZoneRegisteredData(federationContextId, zoneId)
		if err == nil {
			log.Printf("SubscribeZoneController - Zone is already subscribed for federation %s, zoneId %s: %v", federationContextId, zoneId, err)
			utils.HandleProblem(c, http.StatusInternalServerError, "Zone is already subscribed: "+zoneId)
			return
		}

		acceptedZones = append(acceptedZones, zone)

		// check the resources in the zone
		log.Printf("SubscribeZoneController - Getting resources in zone for federation: %s, zoneId: %s", federationContextId, zoneId)
		_, err = zisc.zoneService.GetResourcesInZone(zoneId)
		if err != nil {
			log.Printf("SubscribeZoneController - Error getting resources in zone for federation %s, zoneId %s: %v", federationContextId, zoneId, err)
			utils.HandleProblem(c, http.StatusInternalServerError, "Error getting resources in zone: "+zoneId)
			return
		}
	}

	// make the response body
	log.Printf("SubscribeZoneController - Creating zone registration response for federation: %s", federationContextId)
	zoneRegistrationResponseData := models.ZoneRegistrationResponseData{}

	// for each zone in the acceptedzones, make the corresponding ZoneRegisteredData
	acceptedZoneResourceInfo := []models.ZoneRegisteredData{}
	for _, zone := range acceptedZones {
		zoneRegisteredData := models.ZoneRegisteredData{
			ZoneId:              zone.ZoneId,
			FederationContextId: federationContextId,
		}

		// save the zone registered data to the database
		log.Printf("SubscribeZoneController - Saving zone registered data to database for federation: %s, zoneId: %s", federationContextId, zone.ZoneId)
		err := zisc.zoneService.SaveZoneRegisteredData(zoneRegisteredData)
		if err != nil {
			log.Printf("SubscribeZoneController - Error saving zone registered data for federation %s, zoneId %s: %v", federationContextId, zone.ZoneId, err)
			utils.HandleProblem(c, http.StatusInternalServerError, "Error saving zone registered data: "+zone.ZoneId)
			return
		}

		acceptedZoneResourceInfo = append(acceptedZoneResourceInfo, zoneRegisteredData)
	}

	zoneRegistrationResponseData.AcceptedZoneResourceInfo = acceptedZoneResourceInfo

	log.Printf("SubscribeZoneController - Zone subscription completed successfully for federation: %s, %d zone(s) registered", federationContextId, len(acceptedZoneResourceInfo))
	// return the accepted zones
	c.JSON(http.StatusOK, zoneRegistrationResponseData)
}

// @Summary Unsubscribe from Zone
// @Description Unregisters the origin operator's subscription to a specific availability zone from the partner operator
// @Tags EWBI - ZonesInfoSync
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param zoneId path string true "Zone identifier"
// @Success 200 {object} map[string]string "message: Zone unregistered successfully"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /ewbi/{federationContextId}/zones/{zoneId} [delete]
func (zisc *ZonesInfoSyncController) UnsubscribeZoneController(c *gin.Context) {
	// get the federationContextId from the path
	federationContextId := c.Param("federationContextId")

	// get the zoneId from the path
	zoneId := c.Param("zoneId")
	log.Printf("UnsubscribeZoneController - Starting zone unsubscription for federation: %s, zoneId: %s", federationContextId, zoneId)

	// get the zone registered data for the given zoneId and federationContextId
	log.Printf("UnsubscribeZoneController - Getting zone registered data for federation: %s, zoneId: %s", federationContextId, zoneId)
	_, err := zisc.zoneService.GetZoneRegisteredData(federationContextId, zoneId)
	if err != nil {
		log.Printf("UnsubscribeZoneController - Error getting zone registered data for federation %s, zoneId %s: %v", federationContextId, zoneId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting zone registered data: "+zoneId)
		return
	}

	// delete the zone registered data for the given zoneId and federationContextId
	log.Printf("UnsubscribeZoneController - Deleting zone registered data for federation: %s, zoneId: %s", federationContextId, zoneId)
	err = zisc.zoneService.DeleteZoneRegisteredData(federationContextId, zoneId)
	if err != nil {
		log.Printf("UnsubscribeZoneController - Error deleting zone registered data for federation %s, zoneId %s: %v", federationContextId, zoneId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error deleting zone registered data: "+zoneId)
		return
	}

	log.Printf("UnsubscribeZoneController - Zone unsubscription completed successfully for federation: %s, zoneId: %s", federationContextId, zoneId)
	c.JSON(http.StatusOK, gin.H{"message": "Zone unregistered successfully"})
}

// @Summary Get Zone Details
// @Description Retrieves detailed information about a specific availability zone from the partner operator
// @Tags EWBI - ZonesInfoSync
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param zoneId path string true "Zone identifier"
// @Success 200 {object} models.ZoneDetails
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /ewbi/{federationContextId}/zones/{zoneId} [get]
func (zisc *ZonesInfoSyncController) GetZoneController(c *gin.Context) {
	// get the zoneId from the path
	zoneId := c.Param("zoneId")
	federationContextId := c.Param("federationContextId")
	log.Printf("GetZoneController - Getting zone details for federation: %s, zoneId: %s", federationContextId, zoneId)

	// get the zone details from the database
	zone, err := zisc.zoneService.GetLocalZoneById(zoneId)
	if err != nil {
		log.Printf("GetZoneController - Error getting zone details for federation %s, zoneId %s: %v", federationContextId, zoneId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting zone details")
		return
	}

	log.Printf("GetZoneController - Successfully retrieved zone details for federation: %s, zoneId: %s", federationContextId, zoneId)
	c.JSON(http.StatusOK, zone)
}

// @Summary Get All Local Zones
// @Description Retrieves a list of all availability zones offered by the partner operator
// @Tags EWBI - ZonesInfoSync
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Success 200 {array} models.ZoneDetails
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /ewbi/{federationContextId}/zones [get]
func (zisc *ZonesInfoSyncController) GetAllLocalZonesController(c *gin.Context) {
	federationContextId := c.Param("federationContextId")
	log.Printf("GetAllLocalZonesController - Getting all local zones for federation: %s", federationContextId)

	// ensure latest zones are up to date
	localZones, err := zisc.zoneService.GetLocalZones()
	if err != nil {
		log.Printf("GetAllLocalZonesController - Error getting local zones for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error getting local zones")
		return
	}

	log.Printf("GetAllLocalZonesController - Successfully retrieved %d local zones for federation: %s", len(localZones), federationContextId)
	c.JSON(http.StatusOK, localZones)
}

// @Summary Post MEH Metrics
// @Description Receives metrics data from the orchestrator's Multi-access Edge Host and publishes to Kafka topic
// @Tags EWBI - ZonesInfoSync
// @Accept json
// @Produce json
// @Param federationContextId path string true "Federation Context ID"
// @Param metricsRequestData body dto.OrchMehMetricsRequestData true "MEH Metrics Request Data"
// @Success 200 {object} map[string]string "msgId: Kafka message ID"
// @Failure 400 {object} models.ProblemDetails "Invalid request body"
// @Failure 500 {object} models.ProblemDetails "Error posting metrics to Kafka"
// @Router /ewbi/{federationContextId}/metrics [post]
func (zisc *ZonesInfoSyncController) PostMetricsController(c *gin.Context) {
	federationContextId := c.Param("federationContextId")
	log.Printf("PostMetricsController - Starting metrics posting for federation: %s", federationContextId)

	// get the request body and decode
	log.Printf("PostMetricsController - Binding request body for federation: %s", federationContextId)
	var metricsRequestData dto.OrchMehMetricsRequestData
	if err := c.ShouldBindJSON(&metricsRequestData); err != nil {
		log.Printf("PostMetricsController - Error binding request body for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// post metrics to kafka topic
	metricsTopic := "federation-meh-metrics"
	metricsMessage := metricsRequestData
	log.Printf("PostMetricsController - Publishing metrics to Kafka topic '%s' for federation: %s", metricsTopic, federationContextId)

	msgId, err := zisc.kafkaClientService.Produce(metricsTopic, metricsMessage)
	if err != nil {
		log.Printf("PostMetricsController - Error publishing metrics to Kafka for federation %s: %v", federationContextId, err)
		utils.HandleProblem(c, http.StatusInternalServerError, "Error posting metrics to kafka: "+err.Error())
		return
	}

	log.Printf("PostMetricsController - Metrics posted successfully for federation: %s, msgId: %s", federationContextId, msgId)
	c.JSON(http.StatusOK, gin.H{"msgId": msgId})
}
