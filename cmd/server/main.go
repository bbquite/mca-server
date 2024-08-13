package main

import (
	"encoding/json"
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
	A string `json:"host"`
}

func initOptions() *Options {
	opt := new(Options)

	err := godotenv.Load()
	if err != nil {
		log.Print(".env file not found")
	}

	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		opt.A = envHOST
	} else {
		flag.StringVar(&opt.A, "a", defHost, "server host")
		flag.Parse()
	}

	jsonOptions, _ := json.Marshal(opt)
	log.Printf("Current options: %s", jsonOptions)

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

	if err := srv.Run(opt.A, handler.InitChiRoutes()); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occured while running http server: %v", err)
		}
	}
}
