package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mankings/mec-federator/docs"
	"github.com/mankings/mec-federator/internal/config"
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

func Init(svcs *Services, mdws *Middlewares) *gin.Engine {
	// start gin with default settings
	router := gin.Default()

	// init routes
	initRoutes(router, svcs, mdws)

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
	initEwbiFederationManagementRoutes(router, svcs, mdws)
	initEwbiZoneInfoSyncRoutes(router, svcs, mdws)
	initEwbiArtefactManagementRoutes(router, svcs, mdws)
	initEwbiApplicationInstanceLifecycleManagementRoutes(router, svcs, mdws)
}
