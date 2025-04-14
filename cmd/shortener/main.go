package main

import (
	"log"
	"net/http"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/handler"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
)

func main() {
	// Get app config
	var config = app.GetServerConfig()
	// Init custom logger
	logger.InitLogger(config.LogLevel)

	// repo choice
	var repo repository.RepoInterface
	// set db repo
	if len(config.DBDSN) > 0 {
		dbrepo, err := repository.GetGormRepo(config.DBDSN)
		if err != nil {
			panic(err)
		}
		repo = dbrepo
		logger.Log.Info("Connect to db")
	} else if len(config.FileStoragePath) > 0 {
		//  set file-mem repo
		repo = repository.GetFileRepository(config.FileStoragePath)
	} else {
		// set  mem repo
		repo = repository.GetMemRepository()
	}

	// Bussiness layer (empty)
	serv := service.GetService(repo, config.BaseURL)
	// Handlers
	handler := handler.GetHandler(serv)

	// Middlewares chain
	middlewares := []func(http.Handler) http.Handler{
		middware.Logger,
		middware.GzipMiddleware,
		middware.AuthMiddleware,
	}

	//	Create Router
	router := handler.InitRouter(middlewares)

	// Run server
	server := new(app.Server)
	err := server.Run(config.ServerAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
