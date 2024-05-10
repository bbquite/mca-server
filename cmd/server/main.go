package main

import (
	"context"
	"github.com/bbquite/mca-server/internal/handlers"
	"github.com/bbquite/mca-server/internal/server"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"log"
)

func main() {

	srv := new(server.Server)

	db, err := storage.NewMemStorage()
	if err != nil {
		log.Fatalf("storage initialization error: %s", err.Error())
	}

	metricService := service.NewMetricService(db)
	ctx := context.WithValue(context.Background(), server.ServiceCtx{}, metricService)

	if err := srv.Run("8080", ctx, handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

	defer srv.Shutdown(context.Background())
}
