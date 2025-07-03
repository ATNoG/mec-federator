package services

import (
	"context"
	"fmt"
	"time"

	"github.com/mankings/mec-federator/internal/config"
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
}

// NewFederationService creates a new instance of the FederationServiceImpl
func NewFederationService() *FederationService {
	return &FederationService{}
}

func (fs *FederationService) getFederationCollection() *mongo.Collection {
	return config.GetMongoDatabase().Collection("federations")
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

// GetFederation retrieves a federation from the database using the federationContextId
func (fs *FederationService) GetFederation(federationContextId string) (models.Federation, error) {
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federationContextId}
	federation, err := utils.FetchEntityFromDatabase[models.Federation](collection, filter)
	if err != nil {
		return models.Federation{}, fmt.Errorf("error fetching federation using federationContextId: %v", err)
	}
	return federation, nil
}

// UpdateFederation updates a federation in the database
func (fs *FederationService) UpdateFederation(federation models.Federation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federation.PartnerOP.FederationContextId}
	update := bson.M{"$set": federation}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// PatchFederation updates a federation in the database using patch parameters, then saves it
func (fs *FederationService) PatchFederation(federationContextId string, patchParams models.FederationPatchParams) error {
	federation, err := fs.GetFederation(federationContextId)
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

// RenewFederation renews the federation in the database
func (fs *FederationService) RenewFederation(federation models.Federation) (models.Federation, error) {
	renewalDate := time.Now().AddDate(0, 6, 0)
	expiryDate := renewalDate.AddDate(0, 6, 0)

	federation.PartnerOP.FederationRenewalDate = renewalDate
	federation.PartnerOP.FederationExpiryDate = expiryDate

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federation.PartnerOP.FederationContextId}
	update := bson.M{"$set": federation}
	_, err := collection.UpdateOne(ctx, filter, update)
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

// ExistsFederationWithAccessToken checks if a federation exists in the database using the accessToken
func (fs *FederationService) ExistsFederationWithAccessToken(accessToken string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := fs.getFederationCollection()
	filter := bson.M{"originOP.accessToken.accessToken": accessToken}
	count, _ := collection.CountDocuments(ctx, filter)
	return count > 0
}

// GetFederationContextId retrieves the federationContextId from the database using the accessToken
func (fs *FederationService) GetFederationContextId(accessToken string) (string, error) {
	collection := fs.getFederationCollection()
	filter := bson.M{"originOP.accessToken.accessToken": accessToken}
	federationContextId, err := utils.FetchEntityFromDatabase[string](collection, filter)
	if err != nil {
		return "", fmt.Errorf("error fetching federationContextId using accessToken: %v", err)
	}
	return federationContextId, nil
}

// GetFederationByAccessToken retrieves the federation from the database using the accessToken
func (fs *FederationService) GetFederationByAccessToken(accessToken string) (models.Federation, error) {
	collection := fs.getFederationCollection()
	filter := bson.M{"originOP.accessToken.accessToken": accessToken}
	federation, err := utils.FetchEntityFromDatabase[models.Federation](collection, filter)
	if err != nil {
		return models.Federation{}, fmt.Errorf("error fetching federation using accessToken: %v", err)
	}
	return federation, nil
}

func (fs *FederationService) GetAccessToken(federationContextId string) (string, error) {
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federationContextId}
	accessToken, err := utils.FetchEntityFromDatabase[string](collection, filter)
	if err != nil {
		return "", fmt.Errorf("error fetching accessToken using federationContextId: %v", err)
	}
	return accessToken, nil
}

// GetFederatorUrl retrieves the federatorUrl from the database using the federationContextId
func (fs *FederationService) GetFederatorUrl(federationContextId string) (string, error) {
	collection := fs.getFederationCollection()
	filter := bson.M{"partnerOP.federationContextId": federationContextId}
	federatorUrl, err := utils.FetchEntityFromDatabase[string](collection, filter)
	if err != nil {
		return "", fmt.Errorf("error fetching federatorUrl using federationContextId: %v", err)
	}
	return federatorUrl, nil
}

// Get Federations retrieves all federations from the database
func (fs *FederationService) GetFederations() ([]models.Federation, error) {
	collection := fs.getFederationCollection()
	filter := bson.M{}
	federations, err := utils.FetchEntitiesFromDatabase[models.Federation](collection, filter)
	if err != nil {
		return []models.Federation{}, fmt.Errorf("error fetching federations: %v", err)
	}
	return federations, nil
}

// GetOriginFederations retrieves the federations in which this federator is the origin OP from the database
func (fs *FederationService) GetOriginFederations() ([]models.Federation, error) {
	collection := fs.getFederationCollection()
	filter := bson.M{"isOriginOP": true}
	federations, err := utils.FetchEntitiesFromDatabase[models.Federation](collection, filter)
	if err != nil {
		return []models.Federation{}, fmt.Errorf("error fetching origin federations: %v", err)
	}
	return federations, nil
}

// GetPartnerFederations retrieves the federations in which this federator is the partner OP from the database
func (fs *FederationService) GetPartnerFederations() ([]models.Federation, error) {
	collection := fs.getFederationCollection()
	filter := bson.M{"isOriginOP": false}
	federations, err := utils.FetchEntitiesFromDatabase[models.Federation](collection, filter)
	if err != nil {
		return []models.Federation{}, fmt.Errorf("error fetching partner federations: %v", err)
	}
	return federations, nil
}
