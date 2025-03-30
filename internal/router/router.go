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
	OrchestratorService *services.OrchestratorService
	ZoneService         *services.ZoneService
	HttpClientService   *services.HttpClientService
}

func Init() *gin.Engine {
	// start gin with default settings
	router := gin.Default()

	// init clients
	mongoClient := config.GetMongoClient()
	httpClient := &http.Client{}

	// init services
	authServ := services.NewAuthService(mongoClient)
	fedServ := services.NewFederationService(mongoClient)
	orchServ := services.NewOrchestratorService()
	zoneServ := services.NewZoneService(mongoClient, orchServ, fedServ)
	httpServ := services.NewHttpClientService(httpClient)

	// bundle all services
	services := &Services{
		AuthService:         authServ,
		FederationService:   fedServ,
		OrchestratorService: orchServ,
		ZoneService:         zoneServ,
		HttpClientService:   httpServ,
	}

	// init auth middleware
	authMiddleware := middleware.AuthMiddleware(services.AuthService)

	// init routes
	initRoutes(router, services, authMiddleware)

	// init swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run the server
	router.Run(":" + config.AppConfig.ApiPort)

	return router
}

func initRoutes(router *gin.Engine, svcs *Services, authMiddleware gin.HandlerFunc) {
	// Auth Routes
	initAuthRoutes(router, svcs, authMiddleware)

	// EWBI Routes
	// initFederationAPIManagementRoutes(router, svcs, authMiddleware)
	initEwbiFederationManagementRoutes(router, svcs, authMiddleware)
	initZoneInfoSyncRoutes(router, svcs, authMiddleware)
	initEwbiArtefactManagementRoutes(router, svcs, authMiddleware)

	// SBI Routes
	// interfaces with orchestrators and gives them orders

	// NBI Routes
	// receives orders from the orchestrators
	initNbiFederationManagementRoutes(router, svcs)
	// initNbiArtefactManagementRoutes(router, svcs)
}
