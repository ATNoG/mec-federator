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
	mongoClient := config.GetMongoClient()
	httpClient := &http.Client{}

	// init services
	kafkaServ := services.NewKafkaService()
	authServ := services.NewAuthService(mongoClient)
	fedServ := services.NewFederationService(mongoClient)
	mecServ := services.NewMecSystemService(mongoClient)
	orchServ := services.NewOrchestratorService(kafkaServ)
	zoneServ := services.NewZoneService(mongoClient, orchServ, fedServ)
	httpServ := services.NewHttpClientService(httpClient)

	// bundle all services
	services := &Services{
		AuthService:         authServ,
		FederationService:   fedServ,
		MecSystemService:    mecServ,
		OrchestratorService: orchServ,
		ZoneService:         zoneServ,

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
	// Auth Routes
	initAuthRoutes(router, svcs)

	// EWBI Routes
	// initFederationAPIManagementRoutes(router, svcs, authMiddleware)
	initEwbiFederationManagementRoutes(router, svcs, mdws)
	initZoneInfoSyncRoutes(router, svcs, mdws)
	initEwbiArtefactManagementRoutes(router, svcs, mdws)
}
