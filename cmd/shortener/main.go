package main

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/handler"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"log"
	"net/http"
)

func main() {
	// Init custom logger
	logger.InitLogger("INFO")
	// Get app config
	var config = app.GetServerConfig()

	// mem storage
	memstorage := app.GetStorage()
	// file storage = mem storage + initial read and save changes to file
	repo := repository.GetFileRepository(memstorage, config.FileStoragePath)
	// Bussiness layer (empty)
	serv := service.GetService(repo)
	// Handlers
	handler := handler.GetHandler(serv, config.BaseURL)
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
