package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(host string, mux *chi.Mux) error {

	s.httpServer = &http.Server{
		Addr:           host,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v\n", err)
	} else {
		log.Printf("gracefully stopped\n")
	}
	s.httpServer.RegisterOnShutdown(cancel)

	return s.httpServer.ListenAndServe()
}
