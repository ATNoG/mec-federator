package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/middleware"
	"github.com/mankings/mec-federator/internal/services"
)

type Services struct {
	AuthService       *services.AuthService
	FederationService *services.FederationService
}

func Init() *gin.Engine {
	// start gin with default settings
	router := gin.Default()

	mongoClient := config.GetMongoClient()

	services := &Services{
		AuthService:       services.NewAuthService(mongoClient),
		FederationService: services.NewFederationService(mongoClient),
	}

	authMiddleware := middleware.AuthMiddleware(services.AuthService)

	// init routes
	initRoutes(router, services, authMiddleware)

	// run the server
	router.Run(":" + config.AppConfig.ApiPort)

	return router
}
