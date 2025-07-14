// server
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/handler"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

// Server struct
type Server struct {
	hTTPServer *http.Server
}

// Run(serverAddress string, router *chi.Mux)
func (s *Server) Run(serverAddress string, router *chi.Mux, enableHTTPS bool) error {

	s.hTTPServer = &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}
	if enableHTTPS {
		return http.Serve(autocert.NewListener(serverAddress), router)
	}
	return s.hTTPServer.ListenAndServe()
}

// ServerRun (config app.ServerConfigStruct)
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
	// Create server
	server := new(Server)

	idleChan := make(chan struct{})
	exitChan := make(chan os.Signal, 10)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-exitChan
		if err := server.hTTPServer.Shutdown(context.Background()); err != nil {
			logger.Log.Error("Server shuntdown: %v", zap.Error(err))
		}
		close(idleChan)
	}()

	// Run server
	err := server.Run(config.ServerAddress, router, config.EnableHTTPS)
	if err != nil {
		log.Fatal(err)
	}
	<-idleChan
	repo.Close()
	logger.Log.Info("Server Shutdown gracefully")
}
