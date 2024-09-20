package app

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
	defDatabase        string = ""
)

type serverOptions struct {
	Host            string `json:"HOST"`
	StoreInterval   int64  `json:"STORE_INTERVAL"`
	FilrStoragePath string `json:"FILE_STORAGE_PATH"`
	Restore         bool   `json:"RESTORE"`
	DatabaseDSN     string `json:"DATABASE_DSN"`
}

func initServerOptions(logger *zap.SugaredLogger) *serverOptions {
	opt := new(serverOptions)

	err := godotenv.Load()
	if err != nil {
		logger.Info(".env file not found")
	}

	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		opt.Host = envHOST
	} else {
		flag.StringVar(&opt.Host, "a", defHost, "HOST")
	}

	if envSTOREINTERVAL, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeInterval, err := strconv.ParseInt(envSTOREINTERVAL, 10, 64)
		if err != nil {
			flag.Int64Var(&opt.StoreInterval, "i", defStoreInterval, "STORE_INTERVAL")
		} else {
			opt.StoreInterval = storeInterval
		}
	} else {
		flag.Int64Var(&opt.StoreInterval, "i", defStoreInterval, "STORE_INTERVAL")
	}

	if envFILESTORAGEPATH, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		opt.FilrStoragePath = envFILESTORAGEPATH
	} else {
		flag.StringVar(&opt.FilrStoragePath, "f", defFileStoragePath, "FILE_STORAGE_PATH")
	}

	if envRESTORE, ok := os.LookupEnv("RESTORE"); ok {
		boolValue, err := strconv.ParseBool(envRESTORE)
		if err != nil {
			flag.BoolVar(&opt.Restore, "i", defRestore, "RESTORE")
		}
		opt.Restore = boolValue
	} else {
		flag.BoolVar(&opt.Restore, "r", defRestore, "RESTORE")
	}

	if envDATABASE, ok := os.LookupEnv("DATABASE_DSN"); ok {
		opt.DatabaseDSN = envDATABASE
	} else {
		flag.StringVar(&opt.DatabaseDSN, "d", defDatabase, "DATABASE_DSN")
	}

	flag.Parse()

	jsonOptions, _ := json.Marshal(opt)
	logger.Infof("Server run with options: %s", jsonOptions)

	return opt
}

func RunServer() {

	// ctx := context.Background()

	serverLogger, err := InitLogger()
	if err != nil {
		serverLogger.Fatalf("server logger init error: %v", err)
	}

	opt := initServerOptions(serverLogger)

	var syncSaving = false
	if opt.StoreInterval == 0 {
		syncSaving = true
	}

	db := storage.NewMemStorage()
	serv := service.NewMetricService(db, syncSaving, opt.FilrStoragePath)

	handler, err := handlers.NewHandler(serv, serverLogger)
	if err != nil {
		log.Fatalf("handler construction error: %v", err)
	}

	srv := new(server.Server)
	srv.Run(opt.Host, opt.StoreInterval, opt.FilrStoragePath, opt.Restore, handler.InitChiRoutes(), serv, serverLogger)
}
