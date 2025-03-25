// @title MEC Federator API
// @version 1.0
// @description This is the API documentation for the MEC Federator.
// @host localhost:8000
// @BasePath

// @securityDefinitions.oauth2.clientCredentials
// @tokenUrl http://keycloak.local/realms/federation/protocol/openid-connect/token
// @scope.admin Grants full access
package main

import (
	"log"

	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/router"
)

var err error

func init() {
	err = config.InitAppConfig()
	if err != nil {
		panic(err)
	}

	err = config.InitMongoDB()
	if err != nil {
		panic(err)
	}

	err = config.InitOrchestrators()
	if err != nil {
		panic(err)
	}

	log.Print("Initialized all configurations successfully.")
}

func main() {
	// Initialize the router
	router.Init()
}
