package app

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
	defServerHost     string = "localhost:8080"
	defReportInterval int    = 10 // частота отправки метрик
	defPollInterval   int    = 2  // частота опроса метрик
)

type agentOptions struct {
	Host           string `json:"host"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
}

func initAgentOptions(logger *zap.SugaredLogger) *agentOptions {
	opt := new(agentOptions)

	err := godotenv.Load()
	if err != nil {
		logger.Info(".env file not found")
	}

	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		opt.Host = envHOST
	} else {
		flag.StringVar(&opt.Host, "a", defServerHost, "server host")
	}

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		opt.ReportInterval, _ = strconv.Atoi(envReportInterval)
	} else {
		flag.IntVar(&opt.ReportInterval, "r", defReportInterval, "reportInterval")
	}

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		opt.PollInterval, _ = strconv.Atoi(envPollInterval)
	} else {
		flag.IntVar(&opt.PollInterval, "p", defPollInterval, "pollInterval")
	}

	flag.Parse()

	jsonOptions, _ := json.Marshal(opt)
	logger.Info("Current options: %s", jsonOptions)

	return opt
}

func RunAgent() error {

	agentLogger, err := InitLogger()
	if err != nil {
		log.Fatalf("agent logger init error: %v", err)
	}

	opt := initAgentOptions(agentLogger)

	db := storage.NewMemStorage()
	agentServices := service.NewMetricService(db, false, false, "")
	memStat := new(runtime.MemStats)

	pollTicker := time.NewTicker(time.Duration(opt.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(opt.ReportInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			handlers.MetricsCollect(memStat, agentServices, agentLogger)

		case <-reportTicker.C:
			//err := handlers.MetricsURIRequest(agentServices, opt.host)
			// err := handlers.MetricsJSONRequest(agentServices, opt.Host, agentLogger)
			err := handlers.MetricsMapJSONRequest(agentServices, opt.Host, agentLogger)
			if err != nil {
				agentLogger.Errorf("Falied to make request: \n%v", err)
			}
		}
	}
}
