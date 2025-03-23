package main

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/handler"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"log"
	"net/http"
)

func main() {
	// Init custom logger
	middware.InitLogger("INFO")
	// Get app config
	var config = app.GetServerConfig()

	memstorage := app.GetStorage()
	repo := repository.GetRepository(memstorage)
	serv := service.GetService(repo)
	handler := handler.GetHandler(serv, config.BaseURL)

	// Middwares chain
	middlewares := []func(http.Handler) http.Handler{
		middware.Logger,
		middware.GzipMiddleware,
	}
	router := handler.InitRouter(middlewares)
	server := new(app.Server)
	err := server.Run(config.ServerAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
