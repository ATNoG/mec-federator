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
	clientOptions := options.Client().ApplyURI(uri)
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_USERNAME")

	clientOptions.Auth = &options.Credential{
		Username: username,
		Password: password,
	}

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		logger.Error("Error connecting to MongoDB: ", err)
		return nil, err
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logger.Error("Disconected from MongoDB (deferred)", err)
			panic(err)
		}
	}()

	return client, nil
}
