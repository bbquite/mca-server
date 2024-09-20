package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/bbquite/mca-server/internal/handlers"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	defHost           string = "localhost:8080"
	defReportInterval int    = 10 // частота отправки метрик
	defPollInterval   int    = 2  // частота опроса метрик
)

type Options struct {
	A string `json:"host"`
	R int    `json:"report_interval"`
	P int    `json:"poll_interval"`
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
	}

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		opt.R, _ = strconv.Atoi(envReportInterval)
	} else {
		flag.IntVar(&opt.R, "r", defReportInterval, "reportInterval")
	}

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		opt.P, _ = strconv.Atoi(envPollInterval)
	} else {
		flag.IntVar(&opt.P, "p", defPollInterval, "pollInterval")
	}

	flag.Parse()

	jsonOptions, _ := json.Marshal(opt)
	log.Printf("Current options: %s", jsonOptions)

	return opt
}

func agentRun(opt *Options, logger *zap.SugaredLogger) error {
	db := storage.NewMemStorage()
	agentServices := service.NewMetricService(db, false, "")
	memStat := new(runtime.MemStats)

	pollTicker := time.NewTicker(time.Duration(opt.P) * time.Second)
	reportTicker := time.NewTicker(time.Duration(opt.R) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			handlers.MetricsCollect(memStat, agentServices, logger)
		case <-reportTicker.C:
			//err := handlers.MetricsURIRequest(agentServices, opt.a)
			err := handlers.MetricsJSONRequest(agentServices, opt.A, logger)
			if err != nil {
				logger.Errorf("Falied to make request: \n%v", err)
			}
		}
	}
}

func main() {
	agentLogger, err := initLogger()
	if err != nil {
		log.Fatalf("agent logger init error: %v", err)
	}

	opt := initOptions()

	if err := agentRun(opt, agentLogger); err != nil {
		log.Fatal(err)
	}
}
