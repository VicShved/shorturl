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
	// Init custom logger
	logger.InitLogger("INFO")
	// Get app config
	var config = app.GetServerConfig()

	// repo choice
	var repo repository.RepoInterface
	// set db repo
	if len(config.DBDSN) > 0 {
		// postgres driver
		// pgdriver, err := sql.Open("pgx", config.DBDSN)
		// if err != nil {
		// 	panic(err)
		// }
		// defer pgdriver.Close()
		// dbrepo, err := repository.GetDBRepository(pgdriver)
		// if err != nil {
		// 	panic(err)
		// }
		dbrepo := repository.GetDBRepo(config.DBDSN)
		repo = dbrepo
		logger.Log.Info("Connect to db")
	} else if len(config.FileStoragePath) > 0 {
		//  set file-mem repo
		memstorage := app.GetStorage()
		repo = repository.GetFileRepository(memstorage, config.FileStoragePath)
	} else {
		// set  mem repo
		memstorage := app.GetStorage()
		repo = repository.GetMemRepository(memstorage)
	}

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
