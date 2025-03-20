package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
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
