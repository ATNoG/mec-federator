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
	log.Print("Starting MEC Federator, version 1.8")

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

	// init middlewares
	authMiddleware := middleware.AuthMiddleware(authServ)
	federationExistsMiddleware := middleware.FederationExistsMiddleware(fedServ)

	// bundle middlewares
	middlewares := &router.Middlewares{
		AuthMiddleware:             &authMiddleware,
		FederationExistsMiddleware: &federationExistsMiddleware,
	}

	// start kafka consumers with callbacks
	kafkaServ.StartConsumer(context.Background(), "responses", nil, false) // use default callback for responses, might change in the future

	infrastructureInfoCallback := callbacks.NewInfrastructureInfoCallback(services)
	kafkaServ.StartConsumer(context.Background(), "infrastructure-info", infrastructureInfoCallback.HandleMessage, false)

	federationInfrastructureInfoCallback := callbacks.NewFederationInfrastructureInfoCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation-infrastructure-info", federationInfrastructureInfoCallback.HandleMessage, false)

	newFederationCallback := callbacks.NewNewFederationCallback(services)
	kafkaServ.StartConsumer(context.Background(), "new_federation", newFederationCallback.HandleMessage, true)

	removeFederationCallback := callbacks.NewRemoveFederationCallback(services)
	kafkaServ.StartConsumer(context.Background(), "remove_federation", removeFederationCallback.HandleMessage, true)

	newFederationArtefactCallback := callbacks.NewFederationArtefactNewCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_new_artefact", newFederationArtefactCallback.HandleMessage, true)

	removeFederationArtefactCallback := callbacks.NewFederationArtefactRemoveCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_remove_artefact", removeFederationArtefactCallback.HandleMessage, true)

	newFederationAppiCallback := callbacks.NewFederationAppiNewCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_new_appi", newFederationAppiCallback.HandleMessage, true)

	removeFederationAppiCallback := callbacks.NewFederationRemoveAppiCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_remove_appi", removeFederationAppiCallback.HandleMessage, true)

	enableAppInstKDUCallback := callbacks.NewEnableAppInstanceKDUCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_enable_kdu", enableAppInstKDUCallback.HandleMessage, true)

	disableAppInstKDUCallback := callbacks.NewDisableAppInstanceKDUCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_disable_kdu", disableAppInstKDUCallback.HandleMessage, true)

	migrateAppInstNodeCallback := callbacks.NewFederationMigrateNodeCallback(services)
	kafkaServ.StartConsumer(context.Background(), "federation_migrate_node", migrateAppInstNodeCallback.HandleMessage, true)

	// Initialize the scheduler
	sched := scheduler.NewScheduler(services)
	scheduler.CreateTasks(sched)
	sched.Start()

	// Initialize the router
	router.Init(services, middlewares)

	log.Print("MEC Federator 1.8 started successfully")

}
