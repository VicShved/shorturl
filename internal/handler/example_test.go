package handler_test

import (
	"fmt"
	"net/http"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/handler"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
)

func Example() {
	repo := repository.GetFileRepository(app.ServerConfig.FileStoragePath)
	serv := service.GetService(repo, "")
	handler := handler.GetHandler(serv)
	// Middlewares chain
	middlewares := []func(http.Handler) http.Handler{
		middware.AuthMiddleware,
		middware.Logger,
		middware.GzipMiddleware,
	}

	router := handler.InitRouter(middlewares)
	fmt.Print(router)

}
