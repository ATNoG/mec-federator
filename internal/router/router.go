package router

import (
	"github.com/gin-gonic/gin"
)

func Init() {
	// start gin with default settings
	router := gin.Default()

	// init routes
	initRoutes(router)

	// run the server
	router.Run(":8080")
}
