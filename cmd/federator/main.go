package main

import (
	"github.com/mankings/mec-federator/config"
	"github.com/mankings/mec-federator/router"
)

var (
	logger *config.Logger
)

func main() {
	logger = config.GetLogger("main")
	logger.Info("MEC Federator is starting...")

	err := config.Init()
	if err != nil {
		logger.Error("MEC Federator aborted when initializing configurations!")
		return
	}

	// Initialize the router
	router.Init()
}
