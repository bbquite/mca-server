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

	db := storage.NewMemStorage()
	serv := service.NewMetricService(db)
	handler := handlers.NewHandler(serv)

	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

	defer srv.Shutdown(context.Background())
}
