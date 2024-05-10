package server

import (
	"context"
	"net"
	"net/http"
	"time"
)

type ServiceCtx struct {
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, ctx context.Context, mux *http.ServeMux) error {

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		ConnContext: func(_ context.Context, _ net.Conn) context.Context {
			return ctx
		},
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
