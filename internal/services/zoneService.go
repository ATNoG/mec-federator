package services

import (
	"context"
	"slices"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

/*
 * ZoneService
 *	responsible for managing zones
 */

type ZoneServiceInterface interface {
	GetAllZones() ([]models.ZoneDetails, error)
	GetAllLocalZones() ([]models.ZoneDetails, error)
	GetSubscribedZones(federationContextId string) ([]models.ZoneDetails, error)
}

type ZoneService struct {
	orchestratorService *OrchestratorService
	kafkaClientService  *KafkaClientService
}

func NewZoneService(orchestratorService *OrchestratorService, federationService *FederationService, kafkaClientService *KafkaClientService) *ZoneService {
	return &ZoneService{
		orchestratorService: orchestratorService,
		kafkaClientService:  kafkaClientService,
	}
}

func (z *ZoneService) getZoneDetailsCollection() *mongo.Collection {
	return config.GetMongoDatabase().Collection("zoneDetails")
}

// Updates the database with the local available zones
func (z *ZoneService) UpdateLocalZones(zones []models.ZoneDetails) error {
	// get the latest zones from the orchestrator
	availableZones, err := z.GetAllLocalZones()
	if err != nil {
		return err
	}

	// get the local zones from the database
	localZones, err := z.GetAllLocalZones()
	if err != nil {
		return err
	}

	// compare the available zones with the local zones
	for _, zone := range availableZones {
		if !slices.Contains(localZones, zone) {
			// add the zone to the database
			_, err = z.getZoneDetailsCollection().InsertOne(context.Background(), zone)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Returns all the local zones that are registered for federation
func (z *ZoneService) GetAllLocalZones() ([]models.ZoneDetails, error) {

	// get the local zones from the database
	localZones, err := z.getZoneDetailsCollection().Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	// convert the cursor to a slice of models.ZoneDetails
	var zones []models.ZoneDetails
	for localZones.Next(context.Background()) {
		var zone models.ZoneDetails
		err = localZones.Decode(&zone)
		if err != nil {
			return nil, err
		}

		zones = append(zones, zone)
	}

	return zones, nil
}

// Returns all the zones that are registered for federation
func (z *ZoneService) GetAllZones() ([]models.ZoneDetails, error) {
	return nil, nil
}

// Returns all the zones that are subscribed to a given federation context
func (z *ZoneService) GetSubscribedZones(federationContextId string) ([]models.ZoneDetails, error) {
	// Implementation to get all subscribed zones
	return nil, nil
}
