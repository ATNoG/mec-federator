package services

import (
	"context"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

/*
 * ZoneService
 *	responsible for managing zones
 */

type ZoneServiceInterface interface {
	GetAllZones() ([]models.ZoneDetails, error)
	GetLocalZones() ([]models.ZoneDetails, error)
	GetSubscribedZones(federationContextId string) ([]models.ZoneDetails, error)
}

type ZoneService struct {
	kafkaClientService *KafkaClientService
}

func NewZoneService(kafkaClientService *KafkaClientService) *ZoneService {
	return &ZoneService{
		kafkaClientService: kafkaClientService,
	}
}

func (z *ZoneService) getZoneDetailsCollection() *mongo.Collection {
	return config.GetMongoDatabase().Collection("zoneDetails")
}

// Updates the database with the local available zones
func (z *ZoneService) UpdateLocalZones(zones []models.ZoneDetails) error {
	// get the local zones from the database
	localZones, err := z.GetLocalZones()
	if err != nil {
		return err
	}

	// compare the available zones with the local zones
	for _, zone := range zones {
		// Check if a zone with the same ZoneId and VimId already exists
		zoneExists := false
		for _, localZone := range localZones {
			if localZone.ZoneId == zone.ZoneId && localZone.VimId == zone.VimId {
				zoneExists = true
				break
			}
		}

		if !zoneExists {
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
func (z *ZoneService) GetLocalZones() ([]models.ZoneDetails, error) {

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

// Returns a zone by its id
func (z *ZoneService) GetLocalZoneById(zoneId string) (models.ZoneDetails, error) {
	collection := z.getZoneDetailsCollection()
	filter := bson.M{"zoneId": zoneId}
	zone, err := utils.FetchEntityFromDatabase[models.ZoneDetails](collection, filter)
	return zone, err
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

// Get the vim id of a zone
func (z *ZoneService) GetVimId(zoneId string) (string, error) {
	collection := z.getZoneDetailsCollection()
	filter := bson.M{"zoneId": zoneId}
	zone, err := utils.FetchEntityFromDatabase[models.ZoneDetails](collection, filter)
	return zone.VimId, err
}

// Returns the zone for a given vim id
func (z *ZoneService) GetZoneFromVimId(vimId string) (models.ZoneDetails, error) {
	collection := z.getZoneDetailsCollection()
	filter := bson.M{"vimId": vimId}
	zone, err := utils.FetchEntityFromDatabase[models.ZoneDetails](collection, filter)
	return zone, err
}
