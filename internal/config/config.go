package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiPort               string
	BaseUrl               string
	DbUsername            string
	DbPassword            string
	DbHost                string
	DbPort                string
	OAuth2ClientId        string
	OAuth2Secret          string
	KeycloakTokenEndpoint string
}

var AppConfig *Config

func InitAppConfig() error {
	// Initialize environment variables from .env file
	if err := godotenv.Load(); err == nil {
		log.Print("Warning! Using environment variables from .env file")
	}

	// Initialize AppConfig
	AppConfig = &Config{
		ApiPort:               os.Getenv("API_PORT"),
		BaseUrl:               os.Getenv("BASE_URL"),
		DbUsername:            os.Getenv("MONGO_USERNAME"),
		DbPassword:            os.Getenv("MONGO_PASSWORD"),
		DbHost:                os.Getenv("MONGO_HOST"),
		DbPort:                os.Getenv("MONGO_PORT"),
		OAuth2ClientId:        os.Getenv("OAUTH2_CLIENT_ID"),
		OAuth2Secret:          os.Getenv("OAUTH2_SECRET"),
		KeycloakTokenEndpoint: os.Getenv("OAUTH2_TOKEN_ENDPOINT"),
	}

	slog.Info("Environment Variables", "AppConfig", AppConfig)

	return nil
}
