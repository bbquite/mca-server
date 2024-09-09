package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bbquite/mca-server/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(host string, storeInterval int64, filePath string, restore bool, mux *chi.Mux, service *service.MetricService, logger *zap.SugaredLogger) error {

	s.httpServer = &http.Server{
		Addr:           host,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error occured while running http server: %v", err)
		}
	}()

	if storeInterval > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(storeInterval) * time.Second)
				logger.Debugf("Export storage to %s", filePath)
				err := service.SaveToFile(filePath)
				if err != nil {
					logger.Errorf("error occured while export storage: %v", err)
				}
			}
		}()
	}

	if restore {
		logger.Debugf("Import storage from %s", filePath)
		err := service.LoadFromFile(filePath)
		if err != nil {
			logger.Errorf("error occured while import storage: %v", err)
		}
	}

	sig := <-signalCh
	logger.Info("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}

	logger.Debugf("Export storage to %s", filePath)
	err := service.SaveToFile(filePath)
	if err != nil {
		logger.Errorf("error occured while export storage: %v", err)
	}

	logger.Info("Server shutdown gracefully")

	return nil
}
