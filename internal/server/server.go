package server

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(host string, mux *chi.Mux) error {

	log.Printf("INFO | Server HOST: %s", host)

	s.httpServer = &http.Server{
		Addr:           host,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}
