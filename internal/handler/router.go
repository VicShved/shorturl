package handler

import (
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/go-chi/chi/v5"
)

func GetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", logger.AddLogging(HandlePOST))
	router.Post("/api/shorten", logger.AddLogging(HandlePostJSON))
	router.Get("/{key}", logger.AddLogging(HandleGET))
	return router
}
