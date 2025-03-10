package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mankings/mec-federator/internal/config"
	"github.com/mankings/mec-federator/internal/services"
)

type Services struct {
	AuthService *services.AuthServiceImpl
}

func Init() *gin.Engine {
	// start gin with default settings
	router := gin.Default()

	mongoClient := config.GetMongoClient()

	services := &Services{
		AuthService: services.NewAuthService(mongoClient),
	}

	// init routes
	initRoutes(router, services)

	// run the server
	router.Run()

	return router
}
