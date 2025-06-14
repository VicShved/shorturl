package server

import (
	"log"
	"net/http"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/handler"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	hTTPServer *http.Server
}

func (s *Server) Run(serverAddress string, router *chi.Mux) error {

	s.hTTPServer = &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}
	return s.hTTPServer.ListenAndServe()

}

func ServerRun(config app.ServerConfigStruct) {
	// repo choice
	var repo repository.RepoInterface
	// set db repo
	if len(config.DBDSN) > 0 {
		dbrepo, err := repository.GetGormRepo(config.DBDSN)
		if err != nil {
			panic(err)
		}
		repo = dbrepo
		logger.Log.Info("Connect to db", zap.String("DSN", config.DBDSN))
	} else if len(config.FileStoragePath) > 0 {
		//  set file-mem repo
		repo = repository.GetFileRepository(config.FileStoragePath)
	} else {
		// set  mem repo
		repo = repository.GetMemRepository()
	}

	// Bussiness layer (empty)
	serv := service.GetService(repo, config.BaseURL)
	// Handlers
	handler := handler.GetHandler(serv)

	// Middlewares chain
	middlewares := []func(http.Handler) http.Handler{
		middware.AuthMiddleware,
		middware.Logger,
		middware.GzipMiddleware,
	}

	//	Create Router
	router := handler.InitRouter(middlewares)

	// Run server
	server := new(Server)
	err := server.Run(config.ServerAddress, router)
	if err != nil {
		log.Fatal(err)
	}
}
