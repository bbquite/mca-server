package main

import (
	"context"
	"flag"
	"github.com/bbquite/mca-server/internal/handlers"
	"github.com/bbquite/mca-server/internal/server"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"log"
)

type Options struct {
	a string
}

func main() {

	opt := new(Options)
	flag.StringVar(&opt.a, "a", "localhost:8080", "server host")
	flag.Parse()

	srv := new(server.Server)

	db := storage.NewMemStorage()
	serv := service.NewMetricService(db)
	handler := handlers.NewHandler(serv)

	if err := srv.Run(opt.a, handler.InitChiRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

	defer srv.Shutdown(context.Background())
}
