package config

import (
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	client *mongo.Client
	logger *Logger
)

func Init() error {
	var err error
	logger := GetLogger("config")

	if err = godotenv.Load(); err != nil {
		logger.Error("Error loading .env file: ", err)
		return err
	}

	// Initialize MongoDB
	client, err = InitMongoDB()
	if err != nil {
		logger.Error("Error initializing MongoDB: ", err)
		return err
	}

	return nil
}

func GetMongoClient() *mongo.Client {
	return client
}

func GetLogger(p string) *Logger {
	// Initialize Logger
	logger := NewLogger(p)
	return logger
}
