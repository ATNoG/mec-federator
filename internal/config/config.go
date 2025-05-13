package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OperatorId string
	ApiPort    string
	BaseUrl    string

	OAuth2ClientId        string
	OAuth2ClientSecret    string
	KeycloakTokenEndpoint string

	DbUsername string
	DbPassword string
	DbHost     string
	DbPort     string
	Database   string

	KafkaUsername string
	KafkaPassword string
	KafkaHost     string
	KafkaPort     string
}

var AppConfig *Config

func InitAppConfig() error {
	// Initialize environment variables from .env file
	if err := godotenv.Load(); err == nil {
		log.Print("Warning! Using environment variables from .env file")
	}

	// Initialize AppConfig
	AppConfig = &Config{
		OperatorId: os.Getenv("OPERATOR_ID"),
		ApiPort:    os.Getenv("API_PORT"),
		BaseUrl:    os.Getenv("BASE_URL"),

		OAuth2ClientId:        os.Getenv("OAUTH2_CLIENT_ID"),
		OAuth2ClientSecret:    os.Getenv("OAUTH2_CLIENT_SECRET"),
		KeycloakTokenEndpoint: os.Getenv("OAUTH2_TOKEN_ENDPOINT"),

		DbUsername: os.Getenv("MONGO_USERNAME"),
		DbPassword: os.Getenv("MONGO_PASSWORD"),
		DbHost:     os.Getenv("MONGO_HOST"),
		DbPort:     os.Getenv("MONGO_PORT"),
		Database:   os.Getenv("MONGO_DATABASE"),

		KafkaUsername: os.Getenv("KAFKA_USERNAME"),
		KafkaPassword: os.Getenv("KAFKA_PASSWORD"),
		KafkaHost:     os.Getenv("KAFKA_HOST"),
		KafkaPort:     os.Getenv("KAFKA_PORT"),
	}

	slog.Info("Environment Variables\n" +
		"  OperatorId: " + AppConfig.OperatorId + "\n" +
		"  ApiPort: " + AppConfig.ApiPort + "\n" +
		"  BaseUrl: " + AppConfig.BaseUrl + "\n" +
		"  OAuth2ClientId: " + AppConfig.OAuth2ClientId + "\n" +
		"  OAuth2ClientSecret: " + AppConfig.OAuth2ClientSecret + "\n" +
		"  KeycloakTokenEndpoint: " + AppConfig.KeycloakTokenEndpoint + "\n" +
		"  DbUsername: " + AppConfig.DbUsername + "\n" +
		"  DbPassword: " + AppConfig.DbPassword + "\n" +
		"  DbHost: " + AppConfig.DbHost + "\n" +
		"  DbPort: " + AppConfig.DbPort + "\n" +
		"  Database: " + AppConfig.Database + "\n" +
		"  KafkaUsername: " + AppConfig.KafkaUsername + "\n" +
		"  KafkaPassword: " + AppConfig.KafkaPassword + "\n" +
		"  KafkaHost: " + AppConfig.KafkaHost + "\n" +
		"  KafkaPort: " + AppConfig.KafkaPort)

	return nil
}
