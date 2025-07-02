package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/models/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

/*
 * OrchestratorService
 *	responsible for sending messages and orders to the orchestrator
 */

type OrchestratorServiceInterface interface {
}

type OrchestratorService struct {
	kafkaClientService *KafkaClientService
}

func NewOrchestratorService(kafkaClientService *KafkaClientService) *OrchestratorService {
	return &OrchestratorService{
		kafkaClientService: kafkaClientService,
	}
}

func (s *OrchestratorService) getOrchestratorAppPkgsCollection() *mongo.Collection {
	return config.GetOrchestratorMongoDatabase().Collection("app_pkgs")
}

func (s *OrchestratorService) getOrchestratorAppInstancesCollection() *mongo.Collection {
	return config.GetOrchestratorMongoDatabase().Collection("appis")
}

// Onboards an artefact onto the orchestrator
func (s *OrchestratorService) OnboardAppPkg(appPkg dto.NewAppPkg) (string, error) {
	// Insert the app_pkg into the database
	result, err := s.getOrchestratorAppPkgsCollection().InsertOne(context.Background(), appPkg)
	if err != nil {
		return "", err
	}

	// get the app_pkg id as a string
	appPkgId := result.InsertedID.(bson.ObjectID)

	// make a message to send to the kafka topic
	message := dto.NewAppPkgMessage{
		AppPkgId: appPkgId.Hex(),
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaClientService.Produce("new_app_pkg", message)
	if err != nil {
		return "", err
	}

	// wait for a response
	rsp, err := s.kafkaClientService.WaitForResponse(msgId, 10*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return "", err
	}

	// check status of the response
	status := rsp["status"].(float64)
	if status != 201 {
		return "", errors.New("failed to onboard appPkg")
	}

	return appPkgId.Hex(), nil
}

// Remove an appPkg from the orchestrator
func (s *OrchestratorService) RemoveAppPkg(appPkgId string) error {
	// remove the appPkg from the database
	message := dto.DeleteAppPkgMessage{
		AppPkgId: appPkgId,
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaClientService.Produce("delete_app_pkg", message)
	if err != nil {
		return err
	}

	// wait for a response
	rsp, err := s.kafkaClientService.WaitForResponse(msgId, 10*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return err
	}

	// get status field from response
	status := rsp["status"].(float64)

	// if status is not 204, return an error
	if status != 204 {
		return errors.New("failed to remove appPkg from orchestrator")
	}

	// convert the string id to a bson.ObjectID
	appPkgIdBson, err := bson.ObjectIDFromHex(appPkgId)
	if err != nil {
		return err
	}

	// delete the appPkg from the database
	result, err := s.getOrchestratorAppPkgsCollection().DeleteOne(context.Background(), bson.M{"_id": appPkgIdBson})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("appPkg not found")
	}

	return nil
}

// Get an app pkg from the orchestrator db using its mongo object id
func (s *OrchestratorService) GetAppPkg(objectId string) (dto.OrchAppPkg, error) {
	collection := s.getOrchestratorAppPkgsCollection()
	id, err := bson.ObjectIDFromHex(objectId)
	if err != nil {
		return dto.OrchAppPkg{}, err
	}
	filter := bson.M{"_id": id}

	var appPkg dto.OrchAppPkg
	err = collection.FindOne(context.Background(), filter).Decode(&appPkg)
	if err != nil {
		return dto.OrchAppPkg{}, err
	}

	return appPkg, nil
}

// Instantiate an appPkg
func (s *OrchestratorService) InstantiateAppPkg(appPkgId string, vimId string, config string) (string, error) {
	// make a message to send to the kafka topic
	message := dto.InstantiateAppPkgMessage{
		AppPkgId:    appPkgId,
		Name:        "test",
		Description: "test",
		VimId:       vimId,
		Config:      config,
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaClientService.Produce("instantiate_app_pkg", message)
	if err != nil {
		return "", err
	}

	// wait for a response
	rsp, err := s.kafkaClientService.WaitForResponse(msgId, 20*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return "", err
	}

	// get status field from response
	status := rsp["status"].(float64)

	// if status is not 201, return an error
	if status != 201 {
		slog.Warn("failed to instantiate appPkg", "status", status)
		slog.Warn("response", "response", rsp)
		return "", errors.New("failed to instantiate appPkg")
	}

	// get the appInstanceId from the response
	appiId := rsp["appi_id"].(string)

	// return the appInstanceId
	return appiId, nil
}

// Delete an application instance
func (s *OrchestratorService) TerminateAppi(appiId string) error {
	// make a message to send to the kafka topic
	message := dto.TerminateAppiMessage{
		AppiId: appiId,
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaClientService.Produce("terminate_app_pkg", message)
	if err != nil {
		return err
	}

	// wait for a response
	rsp, err := s.kafkaClientService.WaitForResponse(msgId, 10*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return err
	}

	// get status field from response
	status := rsp["status"].(float64)

	// if status is not 204, return an error
	if status < 200 || status >= 300 {
		return errors.New("failed to terminate app instance")
	}

	return nil
}

// Enable a KDU of an app instance
func (s *OrchestratorService) EnableAppInstanceKDU(appdId string, kduId string, nsId string, node string) error {
	// make a message to send to the kafka topic
	message := dto.EnableAppInstanceKDUMessage{
		AppdId: appdId,
		KDUId:  kduId,
		NSId:   nsId,
		Node:   node,
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaClientService.Produce("enable_kdu", message)
	if err != nil {
		return err
	}

	// wait for a response
	rsp, err := s.kafkaClientService.WaitForResponse(msgId, 10*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return err
	}

	// get status field from response
	status := rsp["status"].(string)

	// if status is not success, return an error
	if status != "success" {
		return errors.New("failed to enable app instance KDU")
	}

	return nil
}

// Disable a KDU of an app instance
func (s *OrchestratorService) DisableAppiKDU(appdId string, kduId string, nsId string) error {
	// make a message to send to the kafka topic
	message := dto.DisableAppInstanceKDUMessage{
		AppdId: appdId,
		KDUId:  kduId,
		NSId:   nsId,
	}

	// send the message to the kafka topic
	msgId, err := s.kafkaClientService.Produce("disable_kdu", message)
	if err != nil {
		return err
	}

	// wait for a response
	rsp, err := s.kafkaClientService.WaitForResponse(msgId, 10*time.Second)
	if err != nil {
		slog.Warn("failed to get response from orchestrator", "error", err)
		return err
	}

	// get status field from response
	status := rsp["status"].(string)

	// if status is not success, return an error
	if status != "success" {
		return errors.New("failed to disable app instance KDU")
	}

	return nil
}

// Get app instance details from the orchestrator db
func (s *OrchestratorService) GetAppi(appInstanceId string) (dto.OrchAppI, error) {
	collection := s.getOrchestratorAppInstancesCollection()
	filter := bson.M{"appi_id": appInstanceId}
	var appInstInfo dto.OrchAppI
	err := collection.FindOne(context.Background(), filter).Decode(&appInstInfo)
	if err != nil {
		return dto.OrchAppI{}, err
	}

	return appInstInfo, nil
}

// Get app instances from the orchestrator db from a list of appi ids
func (s *OrchestratorService) GetAppis(appInstanceIds []string) ([]dto.OrchAppI, error) {
	collection := s.getOrchestratorAppInstancesCollection()
	filter := bson.M{"appi_id": bson.M{"$in": appInstanceIds}}
	var appInsts []dto.OrchAppI
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &appInsts)
	if err != nil {
		return nil, err
	}

	return appInsts, nil
}
