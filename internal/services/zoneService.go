package services

import (
	"github.com/mankings/mec-federator/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ZoneServiceInterface interface {
	GetAllZones() ([]models.ZoneDetails, error)
	GetAllLocalZones() ([]models.ZoneDetails, error)
	GetSubscribedZones(federationContextId string) ([]models.ZoneDetails, error)
}

type ZoneService struct {
	mongoClient         *mongo.Client
	orchestratorService *OrchestratorService
	federationService   *FederationService
}

func NewZoneService(mongoClient *mongo.Client, orchestratorService *OrchestratorService, federationService *FederationService) *ZoneService {
	return &ZoneService{
		mongoClient:         mongoClient,
		orchestratorService: orchestratorService,
		federationService:   federationService,
	}
}

func (z *ZoneService) getZoneDetailsCollection() *mongo.Collection {
	return z.mongoClient.Database("mec").Collection("zoneDetails")
}

// Returns all the local zones that are registered for federation
func (z *ZoneService) GetAllLocalZones() ([]models.ZoneDetails, error) {
	// Implementation to get all local zones
	return nil, nil
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
