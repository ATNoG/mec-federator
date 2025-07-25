package config

import (
	"context"
	"log"
	"log/slog"

	"github.com/mankings/mec-federator/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	ctx         context.Context
	mongoClient *mongo.Client
)

func InitMongoDB() error {
	DatabaseURI := "mongodb://" + AppConfig.DbUsername + ":" + AppConfig.DbPassword + "@" + AppConfig.DbHost + ":" + AppConfig.DbPort + "/?directConnection=true"
	log.Printf("Connecting to MongoDB at %s", DatabaseURI)

	mongoConnection := options.Client().ApplyURI(DatabaseURI)

	var err error
	mongoClient, err = mongo.Connect(mongoConnection)
	if err != nil {
		slog.Error("Error connecting to MongoDB", "error", err)
		return err
	}

	ctx = context.TODO()
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		slog.Error("Error pinging MongoDB", "error", err)
		return err
	}

	log.Printf("Successfully connected to MongoDB")
	return nil
}

// returns the MongoDB client
func GetMongoClient() *mongo.Client {
	return mongoClient
}

// returns the MongoDB database
func GetMongoDatabase() *mongo.Database {
	return mongoClient.Database(AppConfig.Database)
}

// returns the Orchestrator MongoDB database
func GetOrchestratorMongoDatabase() *mongo.Database {
	return mongoClient.Database(AppConfig.OrchestratorDatabase)
}

// inits MecSystem information in the database
func InitMecSystemInformation() error {
	collection := mongoClient.Database(AppConfig.Database).Collection("systems")
	orchestratorInfo := models.OrchestratorInfo{
		OperatorId: AppConfig.OperatorId,
		KafkaUrl:   AppConfig.KafkaHost + ":" + AppConfig.KafkaPort,
	}

	_, err := collection.InsertOne(ctx, orchestratorInfo)
	if err != nil {
		return err
	}

	return nil
}
