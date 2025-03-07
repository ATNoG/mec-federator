package config

import (
	"context"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	ctx         context.Context
	mongoClient *mongo.Client
)

func InitMongoDB() error {
	DatabaseURI := "mongodb://" + AppConfig.DbUsername + ":" + AppConfig.DbPassword + "@" + AppConfig.DbHost + ":" + AppConfig.DbPort

	mongoConnection := options.Client().ApplyURI(DatabaseURI)

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
