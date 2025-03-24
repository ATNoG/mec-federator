package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mankings/mec-federator/docs"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/middleware"
	"github.com/mankings/mec-federator/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Services struct {
	FederationHttpClientManager *services.FederationHttpClientManager
	AuthService                 *services.AuthService
	FederationService           *services.FederationService
}

func Init() *gin.Engine {
	// start gin with default settings
	router := gin.Default()

	// init mongo client
	mongoClient := config.GetMongoClient()

	// services to be injected into routes
	services := &Services{
		FederationHttpClientManager: services.NewFederationHttpClientManager(),
		AuthService:                 services.NewAuthService(mongoClient),
		FederationService:           services.NewFederationService(mongoClient),
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
	initTestRoutes(router, authMiddleware)

	// Auth Routes
	initAuthRoutes(router, svcs, authMiddleware)

	// EWBI Routes
	// initFederationAPIManagementRoutes(router, svcs, authMiddleware)
	initEwbiFederationManagementRoutes(router, svcs, authMiddleware)

	// SBI Routes
	// interfaces with orchestrators and gives them orders

	// NBI Routes
	// receives orders from the orchestrators
	initNbiFederationManagementRoutes(router, svcs)
}
