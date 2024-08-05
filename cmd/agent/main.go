package main

import (
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
	defReportInterval int    = 10
	defPollInterval   int    = 2
)

type Options struct {
	a string
	r int
	p int
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
		log.Print("Error loading .env file")
	}

	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		opt.a = envHOST
	} else {
		flag.StringVar(&opt.a, "a", defHost, "server host")
	}

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		opt.r, _ = strconv.Atoi(envReportInterval)
	} else {
		flag.IntVar(&opt.r, "r", defReportInterval, "reportInterval")
	}

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		opt.p, _ = strconv.Atoi(envPollInterval)
	} else {
		flag.IntVar(&opt.p, "p", defPollInterval, "pollInterval")
	}

	flag.Parse()

	return opt
}

func agentRun(opt *Options, logger *zap.SugaredLogger) error {

	logger.Infof("Server HOST: %s", opt.a)

	db := storage.NewMemStorage()
	agentServices := service.NewMetricService(db)
	memStat := new(runtime.MemStats)

	collectTicker := time.NewTicker(time.Duration(opt.p) * time.Second)
	requestTicker := time.NewTicker(time.Duration(opt.p) * time.Second)

	for {
		select {
		case <-collectTicker.C:
			handlers.MetricsCollect(memStat, agentServices, logger)
		case <-requestTicker.C:
			//err := handlers.MetricsURIRequest(agentServices, opt.a)
			err := handlers.MetricsJSONRequest(agentServices, opt.a, logger)
			if err != nil {
				logger.Errorf("Falied to make request: %v", err)
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
