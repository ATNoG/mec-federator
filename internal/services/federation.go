package services

import "go.mongodb.org/mongo-driver/v2/mongo"

/*
 * FederationService
 *	responsible for establishing, maintaining, and terminating a directed federation relationship with a partner OP
 */

type FederationService interface {
}

type FederationServiceImpl struct {
	mongoClient *mongo.Client
}

func NewFederationService(mongoClient *mongo.Client) *FederationServiceImpl {
	return &FederationServiceImpl{
		mongoClient: mongoClient,
	}
}
