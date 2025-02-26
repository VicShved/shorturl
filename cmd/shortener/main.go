package main

import (
	"fmt"
	"github.com/VicShved/shorturl/internal/app"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	router.Post("/", app.HandlePOST)
	router.Get("/{key}", app.HandleGET)

	var config = app.InitServerConfig()

	fmt.Println("Start URL=", config.ServerAddress)
	fmt.Println("Result URL=", config.BaseURL)

	err := http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		panic(err)
	}
}
