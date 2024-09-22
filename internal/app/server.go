package app

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bbquite/mca-server/internal/handlers"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"github.com/bbquite/mca-server/internal/utils"
	"github.com/go-chi/chi/v5"
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

	IsDatabaseUsage bool `json:"DBUsage"`
	IsSyncSaving    bool `json:"SyncSaving"`
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

	opt.IsDatabaseUsage = false
	if opt.DatabaseDSN != "" {
		opt.IsDatabaseUsage = true
	}

	opt.IsSyncSaving = false
	if opt.StoreInterval == 0 && !opt.IsDatabaseUsage {
		opt.IsSyncSaving = true
	}

	jsonOptions, _ := json.Marshal(opt)
	logger.Infof("Server run with options: %s", jsonOptions)

	return opt
}

type server struct {
	httpServer *http.Server
}

func (s *server) runHTTPSever(host string, storeInterval int64, filePath string, restore bool, mux *chi.Mux, service *service.MetricService, logger *zap.SugaredLogger) error {

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

func RunServer() {

	ctx := context.Background()

	serverLogger, err := utils.InitLogger()
	if err != nil {
		serverLogger.Fatalf("server logger init error: %v", err)
	}

	opt := initServerOptions(serverLogger)
	var serv *service.MetricService

	if opt.IsDatabaseUsage {
		storageInstance, err := storage.NewDBStorage(ctx, opt.DatabaseDSN)
		if err != nil {
			log.Fatalf("database connection error: %v", err)
		}
		defer storageInstance.DB.Close()

		err = storageInstance.CheckDatabaseValid()
		if err != nil {
			log.Fatalf("database struct error: %v", err)
		}

		serv = service.NewMetricService(storageInstance, false, true, "")

	} else {
		storageInstance := storage.NewMemStorage()
		serv = service.NewMetricService(storageInstance, opt.IsSyncSaving, false, opt.FilrStoragePath)
	}

	handler, err := handlers.NewHandler(serv, serverLogger)
	if err != nil {
		log.Fatalf("handler construction error: %v", err)
	}

	srv := new(server)
	srv.runHTTPSever(opt.Host, opt.StoreInterval, opt.FilrStoragePath, opt.Restore, handler.InitChiRoutes(), serv, serverLogger)
}
