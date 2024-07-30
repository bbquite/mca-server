package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/bbquite/mca-server/internal/handlers"
	"github.com/bbquite/mca-server/internal/server"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"github.com/joho/godotenv"
)

const (
	defHost string = "localhost:8080"
)

type Options struct {
	a string
}

func initOptions() *Options {
	opt := new(Options)

	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		opt.a = envHOST
	} else {
		flag.StringVar(&opt.a, "a", defHost, "server host")
		flag.Parse()
	}

	return opt
}

func main() {
	opt := initOptions()
	srv := new(server.Server)
	db := storage.NewMemStorage()
	serv := service.NewMetricService(db)

	handler, err := handlers.NewHandler(serv)
	if err != nil {
		log.Fatalf("handler construction error: %v", err)
	}

	if err := srv.Run(opt.a, handler.InitChiRoutes()); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occured while running http server: %v", err)
		}
	}
}
