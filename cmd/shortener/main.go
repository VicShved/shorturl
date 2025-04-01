package main

import (
	"database/sql"
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
	// Init custom logger
	logger.InitLogger("INFO")
	// Get app config
	var config = app.GetServerConfig()

	var repo service.SaverReader
	// set db repo

	if len(config.DBDSN) > 0 {
		// postgres driver
		pgdriver, err := sql.Open("pgx", config.DBDSN)
		if err != nil {
			panic(err)
		}
		defer pgdriver.Close()
		dbrepo, err := repository.GetDBRepository(pgdriver)
		if err != nil {
			panic(err)
		}
		repo = dbrepo
		logger.Log.Info("Connect to db")
	}

	//  mem storage
	memstorage := app.GetStorage()

	if (repo == nil) && len(config.FileStoragePath) > 0 {
		repo = repository.GetFileRepository(memstorage, config.FileStoragePath)
	}

	if repo == nil {
		repo = repository.GetMemRepository(memstorage)
	}

	// file storage = mem storage + initial read and save changes to file

	// Bussiness layer (empty)
	serv := service.GetService(repo, config.BaseURL)
	// Handlers
	handler := handler.GetHandler(serv)
	// Middlewares chain
	middlewares := []func(http.Handler) http.Handler{
		middware.Logger,
		middware.GzipMiddleware,
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
