package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/bbquite/mca-server/internal/handlers"
	"github.com/bbquite/mca-server/internal/server"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	defHost            string = "localhost:8080"
	defStoreInterval   int64  = 300
	defFileStoragePath string = "backup.json"
	defRestore         bool   = true
)

type Options struct {
	A string `json:"HOST"`
	I int64  `json:"STORE_INTERVAL"`
	F string `json:"FILE_STORAGE_PATH"`
	R bool   `json:"RESTORE"`
}

func initLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	defer logger.Sync()

	return sugar, nil
}

func initOptions(logger *zap.SugaredLogger) *Options {
	opt := new(Options)

	err := godotenv.Load()
	if err != nil {
		logger.Info(".env file not found")
	}

	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		opt.A = envHOST
	} else {
		flag.StringVar(&opt.A, "a", defHost, "HOST")

	}

	if envSTOREINTERVAL, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeInterval, err := strconv.ParseInt(envSTOREINTERVAL, 10, 64)
		if err != nil {
			flag.Int64Var(&opt.I, "i", defStoreInterval, "STORE_INTERVAL")
		} else {
			opt.I = storeInterval
		}
	} else {
		flag.Int64Var(&opt.I, "i", defStoreInterval, "STORE_INTERVAL")
	}

	if envFILESTORAGEPATH, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		opt.F = envFILESTORAGEPATH
	} else {
		flag.StringVar(&opt.F, "f", defFileStoragePath, "FILE_STORAGE_PATH")
	}

	if envRESTORE, ok := os.LookupEnv("RESTORE"); ok {
		opt.F = envRESTORE
	} else {
		flag.BoolVar(&opt.R, "r", defRestore, "RESTORE")
	}

	flag.Parse()

	jsonOptions, _ := json.Marshal(opt)
	logger.Infof("Server run with options: %s", jsonOptions)

	return opt
}

func main() {

	serverLogger, err := initLogger()
	if err != nil {
		log.Fatalf("server logger init error: %v", err)
	}

	opt := initOptions(serverLogger)
	db := storage.NewMemStorage()
	serv := service.NewMetricService(db)

	handler, err := handlers.NewHandler(serv, serverLogger)
	if err != nil {
		log.Fatalf("handler construction error: %v", err)
	}

	srv := new(server.Server)
	srv.Run(opt.A, opt.I, opt.F, opt.R, handler.InitChiRoutes(), serv, serverLogger)
}
