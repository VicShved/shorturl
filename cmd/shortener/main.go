package main

import (
	"fmt"
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	logger.Initialize("INFO")
	router := chi.NewRouter()
	router.Post("/", logger.AddLogging(app.HandlePOST))
	router.Get("/{key}", logger.AddLogging(app.HandleGET))

	var config = app.InitServerConfig()

	fmt.Println("Start URL=", config.ServerAddress)
	fmt.Println("Result URL=", config.BaseURL)

	err := http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		panic(err)
	}
}
