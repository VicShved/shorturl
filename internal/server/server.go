// server
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// Init Server(serverAddress string, router *chi.Mux)
func (s *Server) Init(serverAddress string, router *chi.Mux) {

	s.hTTPServer = &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}
}

// Run(serverAddress string, router *chi.Mux)
func (s *Server) Run(enableHTTPS bool) error {
	if enableHTTPS {
		return s.hTTPServer.Serve(autocert.NewListener(s.hTTPServer.Addr))
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
	server.Init(config.ServerAddress, router)

	idleChan := make(chan string)
	exitChan := make(chan os.Signal, 10)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-exitChan
		logger.Log.Info("Catch syscall sygnal")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Shutdown
		if err := server.hTTPServer.Shutdown(ctx); err != nil {
			logger.Log.Error("Server shuntdown: %v", zap.Error(err))
		}
		logger.Log.Info("Send message for shutdown gracefully")
		close(idleChan)
	}()

	// Run server
	err := server.Run(config.EnableHTTPS)
	if err != nil {
		logger.Log.Error("Error", zap.Error(err))
	}

	// Shutdown gracefully
	<-idleChan
	repo.CloseConnect()
	logger.Log.Info("Server Shutdown gracefully")
}
