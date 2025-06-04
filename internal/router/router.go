package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mankings/mec-federator/docs"
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

	HttpClientService *services.HttpClientService
	KafkaService      *services.KafkaService
}

type Middlewares struct {
	AuthMiddleware             *gin.HandlerFunc
	FederationExistsMiddleware *gin.HandlerFunc
}

func Init() *gin.Engine {
	// start gin with default settings
	router := gin.Default()

	// init clients
	httpClient := &http.Client{}

	// init clients
	httpServ := services.NewHttpClientService(httpClient)
	kafkaServ := services.NewKafkaService()

	// init services
	authServ := services.NewAuthService()
	fedServ := services.NewFederationService()
	mecServ := services.NewMecSystemService()
	orchServ := services.NewOrchestratorService(kafkaServ)
	artefactServ := services.NewArtefactService()
	zoneServ := services.NewZoneService(orchServ, fedServ, kafkaServ)

	// bundle all services
	services := &Services{
		AuthService:         authServ,
		FederationService:   fedServ,
		MecSystemService:    mecServ,
		OrchestratorService: orchServ,
		ZoneService:         zoneServ,
		ArtefactService:     artefactServ,

		HttpClientService: httpServ,
		KafkaService:      kafkaServ,
	}

	// init middlewares
	authMiddleware := middleware.AuthMiddleware(services.AuthService)
	federationExistsMiddleware := middleware.FederationExistsMiddleware(services.FederationService)

	// bundle middlewares
	middlewares := &Middlewares{
		AuthMiddleware:             &authMiddleware,
		FederationExistsMiddleware: &federationExistsMiddleware,
	}

	// init routes
	initRoutes(router, services, middlewares)

	// init swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run the server
	router.Run(":" + config.AppConfig.ApiPort)

	return router
}

func initRoutes(router *gin.Engine, svcs *Services, mdws *Middlewares) {
	// Health Routes
	initHealthRoutes(router)

	// Auth Routes
	initAuthRoutes(router, svcs)

	// EWBI Routes
	// initFederationAPIManagementRoutes(router, svcs, authMiddleware)
	initEwbiFederationManagementRoutes(router, svcs, mdws)
	initEwbiZoneInfoSyncRoutes(router, svcs, mdws)
	initEwbiArtefactManagementRoutes(router, svcs, mdws)
	initEwbiApplicationInstanceLifecycleManagementRoutes(router, svcs, mdws)
}
