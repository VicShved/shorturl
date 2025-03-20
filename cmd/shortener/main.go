package main

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/handler"
	"github.com/VicShved/shorturl/internal/logger"
	"log"
)

func main() {
	logger.Initialize("INFO")
	var config = app.InitServerConfig()
	router := handler.GetRouter()
	server := new(app.Server)
	err := server.Run(config.ServerAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
