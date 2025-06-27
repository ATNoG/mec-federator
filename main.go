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
	"context"
	"log"

	"github.com/mankings/mec-federator/internal/callbacks"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/middleware"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/scheduler"
	"github.com/mankings/mec-federator/internal/services"
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
	log.Print("Starting MEC Federator")

	// init services
	httpServ := services.NewHttpClientService()
	kafkaServ := services.NewKafkaClientService()

	authServ := services.NewAuthService()
	mecServ := services.NewMecSystemService()
	fedServ := services.NewFederationService()
	artefactServ := services.NewArtefactService()
	orchServ := services.NewOrchestratorService(kafkaServ)
	appInstanceServ := services.NewAppInstanceService(kafkaServ)
	zoneServ := services.NewZoneService(kafkaServ)

	// init middlewares
	authMiddleware := middleware.AuthMiddleware(authServ)
	federationExistsMiddleware := middleware.FederationExistsMiddleware(fedServ)

	// start kafka consumers with callbacks
	kafkaServ.StartConsumer(context.Background(), "responses", nil) // use default callback for responses, might change in the future
	newFederationCallback := callbacks.NewNewFederationCallback(authServ, httpServ, kafkaServ, fedServ)
	kafkaServ.StartConsumer(context.Background(), "new_federation", newFederationCallback.HandleMessage)
	infrastructureInfoCallback := callbacks.NewInfrastructureInfoCallback(zoneServ)
	kafkaServ.StartConsumer(context.Background(), "infrastructure-info", infrastructureInfoCallback.HandleMessage)
	removeFederationCallback := callbacks.NewRemoveFederationCallback(authServ, httpServ, kafkaServ, fedServ)
	kafkaServ.StartConsumer(context.Background(), "remove_federation", removeFederationCallback.HandleMessage)

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
		AppInstanceService:  appInstanceServ,
	}

	// bundle middlewares
	middlewares := &router.Middlewares{
		AuthMiddleware:             &authMiddleware,
		FederationExistsMiddleware: &federationExistsMiddleware,
	}

	// Initialize the scheduler
	sched := scheduler.NewScheduler(services)
	scheduler.CreateTasks(sched)
	sched.Start()

	// Initialize the router
	router.Init(services, middlewares)
}
