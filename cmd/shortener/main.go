package main

import (
	"fmt"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/server"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	fmt.Print("Build version: ", buildVersion)
	fmt.Print("Build date: ", buildDate)
	fmt.Print("Build commit: ", buildCommit)
	// Get app config
	var config = app.GetServerConfig()
	// Init custom logger
	logger.InitLogger(config.LogLevel)

	server.ServerRun(*config)
}
