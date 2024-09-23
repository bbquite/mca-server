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

type serverConfig struct {
	Host            string `json:"HOST"`
	StoreInterval   int64  `json:"STORE_INTERVAL"`
	FileStoragePath string `json:"FILE_STORAGE_PATH"`
	Restore         bool   `json:"RESTORE"`
	DatabaseDSN     string `json:"DATABASE_DSN"`

	IsDatabaseUsage bool `json:"DBUsage"`
	IsSyncSaving    bool `json:"SyncSaving"`
}

func initServerConfig(logger *zap.SugaredLogger) *serverConfig {
	cfg := new(serverConfig)

	err := godotenv.Load()
	if err != nil {
		logger.Info(".env file not found")
	}

	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Host = envHOST
	} else {
		flag.StringVar(&cfg.Host, "a", defHost, "HOST")
	}

	if envSTOREINTERVAL, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeInterval, err := strconv.ParseInt(envSTOREINTERVAL, 10, 64)
		if err != nil {
			flag.Int64Var(&cfg.StoreInterval, "i", defStoreInterval, "STORE_INTERVAL")
		} else {
			cfg.StoreInterval = storeInterval
		}
	} else {
		flag.Int64Var(&cfg.StoreInterval, "i", defStoreInterval, "STORE_INTERVAL")
	}

	if envFILESTORAGEPATH, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = envFILESTORAGEPATH
	} else {
		flag.StringVar(&cfg.FileStoragePath, "f", defFileStoragePath, "FILE_STORAGE_PATH")
	}

	if envRESTORE, ok := os.LookupEnv("RESTORE"); ok {
		boolValue, err := strconv.ParseBool(envRESTORE)
		if err != nil {
			flag.BoolVar(&cfg.Restore, "r", defRestore, "RESTORE")
		}
		cfg.Restore = boolValue
	} else {
		flag.BoolVar(&cfg.Restore, "r", defRestore, "RESTORE")
	}

	if envDATABASE, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = envDATABASE
	} else {
		flag.StringVar(&cfg.DatabaseDSN, "d", defDatabase, "DATABASE_DSN")
	}

	flag.Parse()

	cfg.IsDatabaseUsage = false
	if cfg.DatabaseDSN != "" {
		cfg.IsDatabaseUsage = true
	}

	cfg.IsSyncSaving = false
	if cfg.StoreInterval == 0 && !cfg.IsDatabaseUsage {
		cfg.IsSyncSaving = true
	}

	jsonConfig, _ := json.Marshal(cfg)
	logger.Infof("Server run with config: %s", jsonConfig)

	return cfg
}

type server struct {
	httpServer *http.Server
}

func (s *server) runHTTPSever(cfg *serverConfig, mux *chi.Mux, service *service.MetricService, logger *zap.SugaredLogger) error {

	s.httpServer = &http.Server{
		Addr:           cfg.Host,
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

	if cfg.StoreInterval > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
				logger.Debugf("Export storage to %s", cfg.FileStoragePath)
				err := service.SaveToFile(cfg.FileStoragePath)
				if err != nil {
					logger.Errorf("error occured while export storage: %v", err)
				}
			}
		}()
	}

	if cfg.Restore {
		logger.Debugf("Import storage from %s", cfg.FileStoragePath)
		err := service.LoadFromFile(cfg.FileStoragePath)
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

	logger.Debugf("Export storage to %s", cfg.FileStoragePath)
	err := service.SaveToFile(cfg.FileStoragePath)
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

	cfg := initServerConfig(serverLogger)
	var serv *service.MetricService

	if cfg.IsDatabaseUsage {
		storageInstance, err := storage.NewDBStorage(ctx, cfg.DatabaseDSN)
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
		serv = service.NewMetricService(storageInstance, cfg.IsSyncSaving, false, cfg.FileStoragePath)
	}

	handler, err := handlers.NewHandler(serv, serverLogger)
	if err != nil {
		log.Fatalf("handler construction error: %v", err)
	}

	srv := new(server)
	srv.runHTTPSever(cfg, handler.InitChiRoutes(), serv, serverLogger)
}
