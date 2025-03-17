package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

/*
 * FederationService
 *	responsible for establishing, maintaining, and terminating a directed federation relationship with a partner OP
 */

type FederationServiceInterface interface {
	SavePendingFederation(federation models.Federation) error
	UpdateFederation(federation models.Federation) error
	DeleteFederation(federation models.Federation) error
	GetFederationContextId(accessToken string) (string, error)
	GetFederationDetails(federationContextId string) (models.Federation, error)
}

type FederationService struct {
	mongoClient *mongo.Client
}

func (fs *FederationService) getFederationCollection() *mongo.Collection {
	return fs.mongoClient.Database("federatorDb").Collection("federations")
}

// NewFederationService creates a new instance of the FederationServiceImpl
func NewFederationService(mongoClient *mongo.Client) *FederationService {
	return &FederationService{
		mongoClient: mongoClient,
	}
}

// SavePendingFederation saves a pending federation to the database
func (fs *FederationService) CreateFederation(federationRequest models.FederationRequestData) (models.Federation, error) {
	federationResponseData := models.FederationResponseData{
		FederationContextId:   uuid.New().String(),
		FederationExpiryDate:  time.Now().AddDate(0, 6, 0),
		FederationRenewalDate: time.Now().AddDate(0, 3, 0),
	}

	healthInfo := models.FederationHealthInfo{}

	federation := models.Federation{
		PartnerOP:     federationResponseData,
		OriginOP:      federationRequest,
		HealthInfo:    healthInfo,
		IsEstablished: true,
		IsOriginOP:    false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	_, err := collection.InsertOne(ctx, federation)
	return federation, err
}

// UpdateFederation updates a federation in the database
func (fs *FederationService) UpdateFederation(federation models.Federation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federation.PartnerOP.FederationContextId}
	update := bson.M{"$set:": federation}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// DeleteFederation deletes a federation from the database using the federationContextId
func (fs *FederationService) DeleteFederation(federationContextId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federationContextId}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}
