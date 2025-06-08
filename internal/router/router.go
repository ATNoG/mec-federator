package router

import (
	"context"

	"github.com/gin-gonic/gin"
	_ "github.com/mankings/mec-federator/docs"
	"github.com/mankings/mec-federator/internal/callbacks"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/middleware"
	"github.com/mankings/mec-federator/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Services struct {
	AuthService         *services.AuthService
	FederationService   *services.FederationService
	MecSystemService    *services.MecSystemService
	OrchestratorService *services.OrchestratorService
	ZoneService         *services.ZoneService
	ArtefactService     *services.ArtefactService
	AppInstanceService  *services.AppInstanceService

	HttpClientService  *services.HttpClientService
	KafkaClientService *services.KafkaClientService
}

type Middlewares struct {
	AuthMiddleware             *gin.HandlerFunc
	FederationExistsMiddleware *gin.HandlerFunc
}

type Callbacks struct {
	ResponseCallback           *callbacks.ResponseCallback
	InfrastructureInfoCallback *callbacks.InfrastructureInfoCallback
	NewFederationCallback      *callbacks.NewFederationCallback
}

func Init() *gin.Engine {
	// start gin with default settings
	router := gin.Default()

	// init services
	httpServ := services.NewHttpClientService()
	kafkaServ := services.NewKafkaClientService()

	authServ := services.NewAuthService()
	fedServ := services.NewFederationService()
	mecServ := services.NewMecSystemService()
	orchServ := services.NewOrchestratorService(kafkaServ)
	artefactServ := services.NewArtefactService()
	appInstanceServ := services.NewAppInstanceService(kafkaServ)
	zoneServ := services.NewZoneService(orchServ, fedServ, kafkaServ)

	// init middlewares
	authMiddleware := middleware.AuthMiddleware(authServ)
	federationExistsMiddleware := middleware.FederationExistsMiddleware(fedServ)

	// init kafka callbacks
	responseCallback := callbacks.NewResponseCallback()
	newFederationCallback := callbacks.NewNewFederationCallback()
	infrastructureInfoCallback := callbacks.NewInfrastructureInfoCallback(zoneServ)

	// start kafka consumers with callbacks
	kafkaServ.StartConsumer(context.Background(), "responses", nil)
	kafkaServ.StartConsumer(context.Background(), "new_federation", newFederationCallback.HandleMessage)
	kafkaServ.StartConsumer(context.Background(), "infrastructure-info", infrastructureInfoCallback.HandleMessage)

	// bundle all services
	services := &Services{
		AuthService:         authServ,
		FederationService:   fedServ,
		MecSystemService:    mecServ,
		OrchestratorService: orchServ,
		ZoneService:         zoneServ,
		ArtefactService:     artefactServ,
		AppInstanceService:  appInstanceServ,

		HttpClientService:  httpServ,
		KafkaClientService: kafkaServ,
	}

	// bundle middlewares
	middlewares := &Middlewares{
		AuthMiddleware:             &authMiddleware,
		FederationExistsMiddleware: &federationExistsMiddleware,
	}

	// bundle callbacks
	callbacks := &Callbacks{
		ResponseCallback:           responseCallback,
		InfrastructureInfoCallback: infrastructureInfoCallback,
		NewFederationCallback:      newFederationCallback,
	}

	// init routes
	initRoutes(router, services, middlewares, callbacks)

	// init swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run the server
	router.Run(":" + config.AppConfig.ApiPort)

	return router
}

func initRoutes(router *gin.Engine, svcs *Services, mdws *Middlewares, cbs *Callbacks) {
	// Health Routes
	initHealthRoutes(router)

	// Auth Routes
	initAuthRoutes(router, svcs)

	// EWBI Routes
	initEwbiFederationManagementRoutes(router, svcs, mdws, cbs)
	initEwbiZoneInfoSyncRoutes(router, svcs, mdws, cbs)
	initEwbiArtefactManagementRoutes(router, svcs, mdws, cbs)
	initEwbiApplicationInstanceLifecycleManagementRoutes(router, svcs, mdws, cbs)
}
