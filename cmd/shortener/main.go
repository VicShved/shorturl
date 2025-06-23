package main

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/server"
)

func main() {
	// Get app config
	var config = app.GetServerConfig()
	// Init custom logger
	logger.InitLogger(config.LogLevel)

	server.ServerRun(*config)
}
