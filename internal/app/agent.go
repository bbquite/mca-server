package app

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/bbquite/mca-server/internal/handlers"
	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"github.com/bbquite/mca-server/internal/utils"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	defServerHost     string = "localhost:8080"
	defReportInterval int    = 10 // частота отправки метрик
	defPollInterval   int    = 2  // частота опроса метрик
)

type agentConfig struct {
	Host           string `json:"host"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
}

func initAgentConfig(logger *zap.SugaredLogger) *agentConfig {
	cfg := new(agentConfig)

	err := godotenv.Load()
	if err != nil {
		logger.Info(".env file not found")
	}

	// Хост
	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Host = envHOST
	} else {
		flag.StringVar(&cfg.Host, "a", defServerHost, "server host")
	}

	// Частота отправки метрик
	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		cfg.ReportInterval, _ = strconv.Atoi(envReportInterval)
	} else {
		flag.IntVar(&cfg.ReportInterval, "r", defReportInterval, "reportInterval")
	}

	// Частота опроса метрик
	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		cfg.PollInterval, _ = strconv.Atoi(envPollInterval)
	} else {
		flag.IntVar(&cfg.PollInterval, "p", defPollInterval, "pollInterval")
	}

	flag.Parse()

	jsonConfig, _ := json.Marshal(cfg)
	logger.Infof("Current agent config: %s", jsonConfig)

	return cfg
}

func collectMetrics(memStat *runtime.MemStats, services *service.MetricService, logger *zap.SugaredLogger) {

	logger.Info("collecting metrics")
	runtime.ReadMemStats(memStat)

	_, err := services.AddGaugeItem("Alloc", model.Gauge(memStat.Alloc))
	if err != nil {
		logger.Errorf("metric saving error: Alloc = %v", memStat.Alloc)
	}

	_, err = services.AddGaugeItem("BuckHashSys", model.Gauge(memStat.BuckHashSys))
	if err != nil {
		logger.Errorf("metric saving error: BuckHashSys = %v", memStat.BuckHashSys)
	}

	_, err = services.AddGaugeItem("Frees", model.Gauge(memStat.Frees))
	if err != nil {
		logger.Errorf("metric saving error: Frees = %v", memStat.Frees)
	}

	_, err = services.AddGaugeItem("GCCPUFraction", model.Gauge(memStat.GCCPUFraction))
	if err != nil {
		logger.Errorf("metric saving error: Alloc = %v", memStat.GCCPUFraction)
	}

	_, err = services.AddGaugeItem("GCSys", model.Gauge(memStat.GCSys))
	if err != nil {
		logger.Errorf("metric saving error: GCSys = %v", memStat.GCSys)
	}

	_, err = services.AddGaugeItem("HeapAlloc", model.Gauge(memStat.HeapAlloc))
	if err != nil {
		logger.Errorf("metric saving error: HeapAlloc = %v", memStat.HeapAlloc)
	}

	_, err = services.AddGaugeItem("HeapIdle", model.Gauge(memStat.HeapIdle))
	if err != nil {
		logger.Errorf("metric saving error: HeapIdle = %v", memStat.HeapIdle)
	}

	_, err = services.AddGaugeItem("HeapInuse", model.Gauge(memStat.HeapInuse))
	if err != nil {
		logger.Errorf("metric saving error: HeapInuse = %v", memStat.HeapInuse)
	}

	_, err = services.AddGaugeItem("HeapObjects", model.Gauge(memStat.HeapObjects))
	if err != nil {
		logger.Errorf("metric saving error: HeapObjects = %v", memStat.HeapObjects)
	}

	_, err = services.AddGaugeItem("HeapReleased", model.Gauge(memStat.HeapReleased))
	if err != nil {
		logger.Errorf("metric saving error: HeapReleased = %v", memStat.HeapReleased)
	}

	_, err = services.AddGaugeItem("HeapSys", model.Gauge(memStat.HeapSys))
	if err != nil {
		logger.Errorf("metric saving error: HeapSys = %v", memStat.HeapSys)
	}

	_, err = services.AddGaugeItem("LastGC", model.Gauge(memStat.LastGC))
	if err != nil {
		logger.Errorf("metric saving error: LastGC = %v", memStat.LastGC)
	}

	_, err = services.AddGaugeItem("Lookups", model.Gauge(memStat.Lookups))
	if err != nil {
		logger.Errorf("metric saving error: Lookups = %v", memStat.Lookups)
	}

	_, err = services.AddGaugeItem("MCacheInuse", model.Gauge(memStat.MCacheInuse))
	if err != nil {
		logger.Errorf("metric saving error: MCacheInuse = %v", memStat.MCacheInuse)
	}

	_, err = services.AddGaugeItem("MCacheSys", model.Gauge(memStat.MCacheSys))
	if err != nil {
		logger.Errorf("metric saving error: MCacheSys = %v", memStat.MCacheSys)
	}

	_, err = services.AddGaugeItem("MSpanInuse", model.Gauge(memStat.MSpanInuse))
	if err != nil {
		logger.Errorf("metric saving error: MSpanInuse = %v", memStat.MSpanInuse)
	}

	_, err = services.AddGaugeItem("MSpanSys", model.Gauge(memStat.MSpanSys))
	if err != nil {
		logger.Errorf("metric saving error: MSpanSys = %v", memStat.MSpanSys)
	}

	_, err = services.AddGaugeItem("Mallocs", model.Gauge(memStat.Mallocs))
	if err != nil {
		logger.Errorf("metric saving error: Mallocs = %v", memStat.Mallocs)
	}

	_, err = services.AddGaugeItem("NextGC", model.Gauge(memStat.NextGC))
	if err != nil {
		logger.Errorf("metric saving error: NextGC = %v", memStat.NextGC)
	}

	_, err = services.AddGaugeItem("NumForcedGC", model.Gauge(memStat.NumForcedGC))
	if err != nil {
		logger.Errorf("metric saving error: NumForcedGC = %v", memStat.NumForcedGC)
	}

	_, err = services.AddGaugeItem("NumGC", model.Gauge(memStat.NumGC))
	if err != nil {
		logger.Errorf("metric saving error: NumGC = %v", memStat.NumGC)
	}

	_, err = services.AddGaugeItem("OtherSys", model.Gauge(memStat.OtherSys))
	if err != nil {
		logger.Errorf("metric saving error: OtherSys = %v", memStat.OtherSys)
	}

	_, err = services.AddGaugeItem("PauseTotalNs", model.Gauge(memStat.PauseTotalNs))
	if err != nil {
		logger.Errorf("metric saving error: PauseTotalNs = %v", memStat.PauseTotalNs)
	}

	_, err = services.AddGaugeItem("StackInuse", model.Gauge(memStat.StackInuse))
	if err != nil {
		logger.Errorf("metric saving error: StackInuse = %v", memStat.StackInuse)
	}

	_, err = services.AddGaugeItem("StackSys", model.Gauge(memStat.StackSys))
	if err != nil {
		logger.Errorf("metric saving error: StackSys = %v", memStat.StackSys)
	}

	_, err = services.AddGaugeItem("Sys", model.Gauge(memStat.Sys))
	if err != nil {
		logger.Errorf("metric saving error: NextGC = %v", memStat.Sys)
	}

	_, err = services.AddGaugeItem("TotalAlloc", model.Gauge(memStat.TotalAlloc))
	if err != nil {
		logger.Errorf("metric saving error: NextGC = %v", memStat.TotalAlloc)
	}

	rndValue := rand.Intn(100)
	_, err = services.AddGaugeItem("RandomValue", model.Gauge(rndValue))
	if err != nil {
		logger.Errorf("metric saving error: RandomValue = %v", rndValue)
	}

	_, err = services.AddCounterItem("PollCount", model.Counter(1))
	if err != nil {
		logger.Errorf("metric saving error: PollCount")
	}
}

func RunAgent() error {

	agentLogger, err := utils.InitLogger()
	if err != nil {
		log.Fatalf("agent logger init error: %v", err)
	}

	cfg := initAgentConfig(agentLogger)

	db := storage.NewMemStorage()
	agentServices := service.NewMetricService(db, false, false, "")

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	memStat := new(runtime.MemStats)
	for {
		select {
		case <-pollTicker.C:
			collectMetrics(memStat, agentServices, agentLogger)

		case <-reportTicker.C:
			// err := handlers.MetricsURIRequest(agentServices, cfg.Host, agentLogger)
			err := handlers.MetricsJSONRequest(agentServices, cfg.Host, agentLogger)
			// err := handlers.MetricsPackJSONRequest(agentServices, cfg.Host, agentLogger)
			if err != nil {
				agentLogger.Errorf("Falied to make request: \n%v", err)
			}
		}
	}
}
