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

func valStrFunc(s string) string {
	if s == "" {
		s = "N/A"
	}
	return s
}

func main() {

	fmt.Println("Build version: ", valStrFunc(buildVersion))
	fmt.Println("Build date: ", valStrFunc(buildDate))
	fmt.Println("Build commit: ", valStrFunc(buildCommit))
	// Get app config
	var config = app.GetServerConfig()
	// Init custom logger
	logger.InitLogger(config.LogLevel)

	server.ServerRun(*config)
}
