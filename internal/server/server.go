// server
package server

import (
	"context"
	"log"
	"net"
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
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Server struct
type Server struct {
	address    string
	hTTPServer *http.Server
	gRPCServer *grpc.Server
	mux        cmux.CMux
}

// Init Server(serverAddress string, router *chi.Mux)
func (s *Server) Init(serverAddress string, router *chi.Mux) {
	s.address = serverAddress
	s.hTTPServer = &http.Server{
		Handler: router,
	}
	s.gRPCServer = grpc.NewServer()
}

// Run(serverAddress string, router *chi.Mux)
func (s *Server) Run(enableHTTPS bool) error {
	// Create the listener
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal(err)
	}

	s.mux = cmux.New(listener)
	grpcListener := s.mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := s.mux.Match(cmux.Any())

	grp := errgroup.Group{}
	// run grpc server
	grp.Go(func() error {
		return s.gRPCServer.Serve(grpcListener)
	})
	// run http server
	grp.Go(func() error {
		if enableHTTPS {
			return s.hTTPServer.Serve(autocert.NewListener(s.hTTPServer.Addr))
		}
		return s.hTTPServer.Serve(httpListener)
		// return s.hTTPServer.ListenAndServe()
	})
	// run cmux multiplexer
	grp.Go(func() error {
		return s.mux.Serve()
	})
	// wait group
	err = grp.Wait()
	return err
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
	handler := handler.GetHandler(serv, config.TrustedSubnet)

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
	exitChan := make(chan os.Signal, 3)
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-exitChan
		logger.Log.Info("Catch syscall sygnal")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown
		grp := errgroup.Group{}
		grp.Go(func() error {
			server.mux.Close()
			return nil
		})
		grp.Go(func() error {
			return server.hTTPServer.Shutdown(ctx)
		})
		grp.Go(func() error {
			server.gRPCServer.GracefulStop()
			return nil
		})

		if err := grp.Wait(); err != nil {
			logger.Log.Error("Servers shuntdown: %v", zap.Error(err))
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
	repo.CloseConn()
	logger.Log.Info("Server Shutdown gracefully")
}
