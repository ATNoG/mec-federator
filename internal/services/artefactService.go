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

type ArtefactServiceInterface interface {
	ReceivePartnerArtefact(artefact models.Artefact) error
	UploadArtefact(artefact models.Artefact) (string, error)
	RegisterArtefact(artefact models.Artefact) error
	GetArtefact(federationContextId string, artefactId string) (models.Artefact, error)
	GetArtefactById(artefactId string) (models.Artefact, error)
	RemoveArtefact(federationContextId string, artefactId string) error
	RegisterFile(file models.File) error
	GetFile(federationContextId string, fileId string) (models.File, error)
	RemoveFile(federationContextId string, fileId string) error
}

type ArtefactService struct {
}

func NewArtefactService() *ArtefactService {
	return &ArtefactService{}
}

func (as *ArtefactService) getArtefactCollection() *mongo.Collection {
	return config.GetMongoDatabase().Collection("artefacts")
}

func (as *ArtefactService) getFileCollection() *mongo.Collection {
	return config.GetMongoDatabase().Collection("files")
}

// Register an artefact locally
func (as *ArtefactService) RegisterArtefact(artefact models.Artefact) error {
	artefact.FederationContextId = "local"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := as.getArtefactCollection()
	_, err := collection.InsertOne(ctx, artefact)	
	return err
}

// Save an artefact to the local database
func (as *ArtefactService) SaveArtefact(artefact models.Artefact) error {
	collection := as.getArtefactCollection()
	_, err := collection.InsertOne(context.Background(), artefact)
	return err
}

// Get an artefact from the local database using the federationContextId and artefactId
func (as *ArtefactService) GetArtefact(federationContextId string, artefactId string) (models.Artefact, error) {
	collection := as.getArtefactCollection()
	filter := bson.M{"federationContextId": federationContextId, "artefactId": artefactId}
	artefact, err := utils.FetchEntityFromDatabase[models.Artefact](collection, filter)
	if err != nil {
		return models.Artefact{}, fmt.Errorf("error, could not find artefact with given federationContextId and artefactId: %v", err)
	}
	return artefact, err
}

// Get an artefact from the local database using the artefactId
func (as *ArtefactService) GetArtefactById(artefactId string) (models.Artefact, error) {
	collection := as.getArtefactCollection()
	filter := bson.M{"artefactId": artefactId}
	artefact, err := utils.FetchEntityFromDatabase[models.Artefact](collection, filter)
	if err != nil {
		return models.Artefact{}, fmt.Errorf("error, could not find artefact with given artefactId: %v", err)
	}
	return artefact, err
}

// Remove an artefact from the local database using the federationContextId and artefactId
func (as *ArtefactService) RemoveArtefact(federationContextId string, artefactId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := as.getArtefactCollection()
	filter := bson.M{"federationContextId": federationContextId, "artefactId": artefactId}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

// Register a file locally
func (as *ArtefactService) RegisterFile(file models.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := as.getFileCollection()
	_, err := collection.InsertOne(ctx, file)
	return err
}

// Get a file from the local database using the federationContextId and fileId
func (as *ArtefactService) GetFile(federationContextId string, fileId string) (models.File, error) {
	collection := as.getFileCollection()
	filter := bson.M{"federationContextId": federationContextId, "fileId": fileId}
	file, err := utils.FetchEntityFromDatabase[models.File](collection, filter)
	if err != nil {
		return models.File{}, fmt.Errorf("error, could not find file with given federationContextId and fileId: %v", err)
	}
	return file, err
}

// Remove a file from the local database using the federationContextId and fileId
func (as *ArtefactService) RemoveFile(federationContextId string, fileId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := as.getFileCollection()
	filter := bson.M{"federationContextId": federationContextId, "fileId": fileId}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

// Receive an artefact from a partner
func (as *ArtefactService) ReceivePartnerArtefact(artefact models.Artefact) error {
	return nil
}

// Receive an artefact from the orchestrator
func (as *ArtefactService) UploadArtefact(artefact models.Artefact) (string, error) {
	return "", nil
}
