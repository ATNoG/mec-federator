// @title MEC Federator API
// @version 1.0
// @description This is the API documentation for the MEC Federator.
// @host localhost:8000
// @BasePath /federation/v1

// @securityDefinitions.oauth2.clientCredentials
// @tokenUrl http://keycloak.local/realms/federation/protocol/openid-connect/token
// @scope.admin Grants full access
package main

import (
	"log"

	"github.com/atnog/mec-federator/internal/callbacks"
	"github.com/atnog/mec-federator/internal/config"
	"github.com/atnog/mec-federator/internal/middleware"
	"github.com/atnog/mec-federator/internal/router"
	"github.com/atnog/mec-federator/internal/scheduler"
	"github.com/atnog/mec-federator/internal/services"
)

var err error

func init() {
	err = config.InitAppConfig()
	if err != nil {
		panic(err)
	}

	err = config.InitKafka()
	if err != nil {
		panic(err)
	}

	err = config.InitMongoDB()
	if err != nil {
		panic(err)
	}

	err = config.InitMecSystemInformation()
	if err != nil {
		panic(err)
	}

	log.Print("Initialized all configurations successfully.")
}

func main() {
	log.Print("Starting MEC Federator, version 1.9 (results mode)")

	// init services
	httpServ := services.NewHttpClientService()
	kafkaServ := services.NewKafkaClientService()

	authServ := services.NewAuthService()
	mecServ := services.NewMecSystemService()
	fedServ := services.NewFederationService()
	artefactServ := services.NewArtefactService()
	applicationServ := services.NewApplicationService(kafkaServ)
	orchServ := services.NewOrchestratorService(kafkaServ)
	appInstanceServ := services.NewAppInstanceService(kafkaServ)
	zoneServ := services.NewZoneService(kafkaServ)

	// bundle all services
	services := &router.Services{
		KafkaClientService: kafkaServ,
		HttpClientService:  httpServ,

		AuthService:         authServ,
		OrchestratorService: orchServ,
		MecSystemService:    mecServ,
		FederationService:   fedServ,
		ZoneService:         zoneServ,
		ArtefactService:     artefactServ,
		ApplicationService:  applicationServ,
		AppInstanceService:  appInstanceServ,
	}

	// init middlewares
	authMiddleware := middleware.AuthMiddleware(authServ)
	federationExistsMiddleware := middleware.FederationExistsMiddleware(fedServ)

	// bundle middlewares
	middlewares := &router.Middlewares{
		AuthMiddleware:             &authMiddleware,
		FederationExistsMiddleware: &federationExistsMiddleware,
	}

	// start callbacks
	callbacks.StartCallbacks(services)

	// start kafka consumers with callbacks
	// Initialize the scheduler
	sched := scheduler.NewScheduler(services)
	scheduler.CreateTasks(sched)
	sched.Start()

	// Initialize the router
	router.Init(services, middlewares)

	log.Print("MEC Federator 1.9 (results mode) started successfully")

}
