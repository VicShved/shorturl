package main

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/handler"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"log"
)

func main() {
	logger.Initialize("INFO")
	var config = app.InitServerConfig()

	memstorage := app.GetStorage()
	repo := repository.GetRepository(memstorage)
	serv := service.GetService(repo)
	handler := handler.GetHandler(serv)
	router := handler.InitRouter()

	//router := handler.GetRouter()

	server := new(app.Server)
	err := server.Run(config.ServerAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
