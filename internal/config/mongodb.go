package config

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB() (*mongo.Client, error) {
	logger := GetLogger("mongodb")

	// Create database and connect
	uri := os.Getenv("MONGO_URI")
	logger.Infof("Connecting to MongoDB at %s", uri)

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		logger.Error("Error connecting to MongoDB: ", err)
		return nil, err
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logger.Error("Disconected from MongoDB (defer).", err)
			panic(err)
		}
	}()

	if err := client.Ping(context.Background(), nil); err != nil {
		logger.Error("Error pinging MongoDB: ", err)
		return nil, err
	}

	logger.Info("Connected to MongoDB successfully.")
	return client, nil
}
