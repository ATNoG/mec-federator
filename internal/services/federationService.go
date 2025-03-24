package services

import (
	"context"
	"fmt"
	"time"

	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/utils"
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
	return fs.mongoClient.Database("federationDb").Collection("federations")
}

// NewFederationService creates a new instance of the FederationServiceImpl
func NewFederationService(mongoClient *mongo.Client) *FederationService {
	return &FederationService{
		mongoClient: mongoClient,
	}
}

// CreateFederation saves a new federation to the database
func (fs *FederationService) CreateFederation(federation models.Federation) (models.Federation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	_, err := collection.InsertOne(ctx, federation)
	return federation, err
}

// DeleteFederation deletes a federation from the database using the federationContextId
func (fs *FederationService) DeleteFederation(federationContextId string) error {
	// remove Federation from the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federationContextId}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

// GetFederationFromContextId retrieves a federation from the database using the federationContextId
func (fs *FederationService) GetFederationFromContextId(federationContextId string) (models.Federation, error) {
	var federation models.Federation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federationContextId}
	err := collection.FindOne(ctx, filter).Decode(&federation)
	return federation, err
}

// ExistsFederationWithContextId checks if a federation exists in the database using the federationContextId
func (fs *FederationService) ExistsFederationWithContextId(federationContextId string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federationContextId}
	count, _ := collection.CountDocuments(ctx, filter)
	return count > 0
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

func (fs *FederationService) PatchFederation(federationContextId string, patchParams models.FederationPatchParams) error {
	federation, err := fs.GetFederationFromContextId(federationContextId)
	if err != nil {
		return err
	}

	switch patchParams.ObjectType {
	case "MOBILE_NETWORK_CODES":
		switch patchParams.OperationType {
		case "ADD_CODES":
			utils.AddMobileCodes(federation.OriginOP.OrigOPMobileNetworkCodes, patchParams.AddMobileNetworkIds)

		case "REMOVE_CODES":
			utils.RemoveMobileCodes(federation.OriginOP.OrigOPMobileNetworkCodes, patchParams.RemoveMobileNetworkIds)

		case "UPDATE_CODES":
			federation.OriginOP.OrigOPMobileNetworkCodes = patchParams.AddMobileNetworkIds

		default:
			return fmt.Errorf("unsupported operation type: %s", patchParams.OperationType)
		}

	case "FIXED_NETWORK_CODES":
		switch patchParams.OperationType {
		case "ADD_CODES":
			federation.OriginOP.OrigOPFixedNetworkCodes = utils.AddFixedCodes(federation.OriginOP.OrigOPFixedNetworkCodes, patchParams.AddFixedNetworkIds)

		case "REMOVE_CODES":
			federation.OriginOP.OrigOPFixedNetworkCodes = utils.RemoveFixedCodes(federation.OriginOP.OrigOPFixedNetworkCodes, patchParams.RemoveFixedNetworkIds)

		case "UPDATE_CODES":
			federation.OriginOP.OrigOPFixedNetworkCodes = patchParams.AddFixedNetworkIds

		default:
			return fmt.Errorf("unsupported operation type: %s", patchParams.OperationType)
		}

	default:
		return fmt.Errorf("unsupported object type: %s", patchParams.ObjectType)
	}

	return fs.UpdateFederation(federation)
}

// Helper functions for
