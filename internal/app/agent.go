package app

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
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
	defAgentKey       string = ""
	defRateLimit      int    = 2
)

var memStatsList = [...]string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
	"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys",
	"LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
	"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc",
}

type agentConfig struct {
	Host           string `json:"host"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	Key            string `json:"KEY"`
	RateLimit      int    `json:"rate_limit"`
}

func initAgentConfigENV(cfg *agentConfig) *agentConfig {
	err := godotenv.Load()
	if err != nil {
		log.Print(".env file not found")
	}

	if envHOST, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Host = envHOST
	}

	if envKEY, ok := os.LookupEnv("KEY"); ok {
		cfg.Key = envKEY
	}

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		cfg.ReportInterval, _ = strconv.Atoi(envReportInterval)
	}

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		cfg.PollInterval, _ = strconv.Atoi(envPollInterval)
	}

	if envRateLimit, ok := os.LookupEnv("RATE_LIMIT"); ok {
		cfg.RateLimit, _ = strconv.Atoi(envRateLimit)
	}

	return cfg
}

func collectMetrics(metricStat *runtime.MemStats, services *service.MetricService, logger *zap.SugaredLogger) {

	logger.Info("collecting default metrics")
	runtime.ReadMemStats(metricStat)

	_, err := services.AddGaugeItem("Alloc", model.Gauge(metricStat.Alloc))
	if err != nil {
		logger.Errorf("metric saving error: Alloc = %v", metricStat.Alloc)
	}

	_, err = services.AddGaugeItem("BuckHashSys", model.Gauge(metricStat.BuckHashSys))
	if err != nil {
		logger.Errorf("metric saving error: BuckHashSys = %v", metricStat.BuckHashSys)
	}

	_, err = services.AddGaugeItem("Frees", model.Gauge(metricStat.Frees))
	if err != nil {
		logger.Errorf("metric saving error: Frees = %v", metricStat.Frees)
	}

	_, err = services.AddGaugeItem("GCCPUFraction", model.Gauge(metricStat.GCCPUFraction))
	if err != nil {
		logger.Errorf("metric saving error: Alloc = %v", metricStat.GCCPUFraction)
	}

	_, err = services.AddGaugeItem("GCSys", model.Gauge(metricStat.GCSys))
	if err != nil {
		logger.Errorf("metric saving error: GCSys = %v", metricStat.GCSys)
	}

	_, err = services.AddGaugeItem("HeapAlloc", model.Gauge(metricStat.HeapAlloc))
	if err != nil {
		logger.Errorf("metric saving error: HeapAlloc = %v", metricStat.HeapAlloc)
	}

	_, err = services.AddGaugeItem("HeapIdle", model.Gauge(metricStat.HeapIdle))
	if err != nil {
		logger.Errorf("metric saving error: HeapIdle = %v", metricStat.HeapIdle)
	}

	_, err = services.AddGaugeItem("HeapInuse", model.Gauge(metricStat.HeapInuse))
	if err != nil {
		logger.Errorf("metric saving error: HeapInuse = %v", metricStat.HeapInuse)
	}

	_, err = services.AddGaugeItem("HeapObjects", model.Gauge(metricStat.HeapObjects))
	if err != nil {
		logger.Errorf("metric saving error: HeapObjects = %v", metricStat.HeapObjects)
	}

	_, err = services.AddGaugeItem("HeapReleased", model.Gauge(metricStat.HeapReleased))
	if err != nil {
		logger.Errorf("metric saving error: HeapReleased = %v", metricStat.HeapReleased)
	}

	_, err = services.AddGaugeItem("HeapSys", model.Gauge(metricStat.HeapSys))
	if err != nil {
		logger.Errorf("metric saving error: HeapSys = %v", metricStat.HeapSys)
	}

	_, err = services.AddGaugeItem("LastGC", model.Gauge(metricStat.LastGC))
	if err != nil {
		logger.Errorf("metric saving error: LastGC = %v", metricStat.LastGC)
	}

	_, err = services.AddGaugeItem("Lookups", model.Gauge(metricStat.Lookups))
	if err != nil {
		logger.Errorf("metric saving error: Lookups = %v", metricStat.Lookups)
	}

	_, err = services.AddGaugeItem("MCacheInuse", model.Gauge(metricStat.MCacheInuse))
	if err != nil {
		logger.Errorf("metric saving error: MCacheInuse = %v", metricStat.MCacheInuse)
	}

	_, err = services.AddGaugeItem("MCacheSys", model.Gauge(metricStat.MCacheSys))
	if err != nil {
		logger.Errorf("metric saving error: MCacheSys = %v", metricStat.MCacheSys)
	}

	_, err = services.AddGaugeItem("MSpanInuse", model.Gauge(metricStat.MSpanInuse))
	if err != nil {
		logger.Errorf("metric saving error: MSpanInuse = %v", metricStat.MSpanInuse)
	}

	_, err = services.AddGaugeItem("MSpanSys", model.Gauge(metricStat.MSpanSys))
	if err != nil {
		logger.Errorf("metric saving error: MSpanSys = %v", metricStat.MSpanSys)
	}

	_, err = services.AddGaugeItem("Mallocs", model.Gauge(metricStat.Mallocs))
	if err != nil {
		logger.Errorf("metric saving error: Mallocs = %v", metricStat.Mallocs)
	}

	_, err = services.AddGaugeItem("NextGC", model.Gauge(metricStat.NextGC))
	if err != nil {
		logger.Errorf("metric saving error: NextGC = %v", metricStat.NextGC)
	}

	_, err = services.AddGaugeItem("NumForcedGC", model.Gauge(metricStat.NumForcedGC))
	if err != nil {
		logger.Errorf("metric saving error: NumForcedGC = %v", metricStat.NumForcedGC)
	}

	_, err = services.AddGaugeItem("NumGC", model.Gauge(metricStat.NumGC))
	if err != nil {
		logger.Errorf("metric saving error: NumGC = %v", metricStat.NumGC)
	}

	_, err = services.AddGaugeItem("OtherSys", model.Gauge(metricStat.OtherSys))
	if err != nil {
		logger.Errorf("metric saving error: OtherSys = %v", metricStat.OtherSys)
	}

	_, err = services.AddGaugeItem("PauseTotalNs", model.Gauge(metricStat.PauseTotalNs))
	if err != nil {
		logger.Errorf("metric saving error: PauseTotalNs = %v", metricStat.PauseTotalNs)
	}

	_, err = services.AddGaugeItem("StackInuse", model.Gauge(metricStat.StackInuse))
	if err != nil {
		logger.Errorf("metric saving error: StackInuse = %v", metricStat.StackInuse)
	}

	_, err = services.AddGaugeItem("StackSys", model.Gauge(metricStat.StackSys))
	if err != nil {
		logger.Errorf("metric saving error: StackSys = %v", metricStat.StackSys)
	}

	_, err = services.AddGaugeItem("Sys", model.Gauge(metricStat.Sys))
	if err != nil {
		logger.Errorf("metric saving error: NextGC = %v", metricStat.Sys)
	}

	_, err = services.AddGaugeItem("TotalAlloc", model.Gauge(metricStat.TotalAlloc))
	if err != nil {
		logger.Errorf("metric saving error: NextGC = %v", metricStat.TotalAlloc)
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

	cfgFlags := new(agentConfig)

	flag.StringVar(&cfgFlags.Host, "a", defServerHost, "server host")
	flag.StringVar(&cfgFlags.Key, "k", defAgentKey, "KEY")
	flag.IntVar(&cfgFlags.ReportInterval, "r", defReportInterval, "reportInterval")
	flag.IntVar(&cfgFlags.PollInterval, "p", defPollInterval, "pollInterval")
	flag.IntVar(&cfgFlags.RateLimit, "l", defRateLimit, "rateLimit")
	flag.Parse()

	cfg := initAgentConfigENV(cfgFlags)

	agentLogger, err := utils.InitLogger()
	if err != nil {
		log.Fatalf("agent logger init error: %v", err)
	}

	jsonConfig, _ := json.Marshal(cfg)
	agentLogger.Infof("Current agent config: %s", jsonConfig)

	db := storage.NewMemStorage()
	agentServices, err := service.NewMetricService(db, false, false, "")
	if err != nil {
		log.Fatalf("service construction error: %v", err)
	}

	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	metricStat := new(runtime.MemStats)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-pollTicker.C:
				collectMetrics(metricStat, agentServices, agentLogger)

			case <-reportTicker.C:
				// err := handlers.SendMetricsURI(agentServices, cfg.Host, agentLogger)
				// err := handlers.SendMetricsJSON(agentServices, cfg.Host, cfg.Key, agentLogger)
				err := handlers.SendMetricsPackJSON(agentServices, cfg.Host, cfg.Key, agentLogger)
				if err != nil {
					agentLogger.Errorf("Falied to make request: \n%v", err)
				}
			}
		}
	}()

	sig := <-signalCh
	agentLogger.Info("Received signal: %v\n", sig)

	pollTicker.Stop()
	reportTicker.Stop()

	agentLogger.Info("Agent shutdown gracefully")
	return nil
}

func collectMetricsCPU(ctx context.Context, wg *sync.WaitGroup, pollInterval int, services *service.MetricService, logger *zap.SugaredLogger) {

	for {
		select {
		case <-ctx.Done():
			logger.Info("--> collect CPU metrics goroutine exit")
			wg.Done()
			return

		case <-time.After(time.Duration(pollInterval) * time.Second):
			logger.Info("collecting CPU metrics")
			memoryStat, _ := mem.VirtualMemory()

			_, err := services.AddGaugeItem("TotalMemory", model.Gauge(memoryStat.Total))
			if err != nil {
				logger.Errorf("metric saving error: TotalMemory = %v", memoryStat.Total)
			}

			_, err = services.AddGaugeItem("FreeMemory", model.Gauge(memoryStat.Free))
			if err != nil {
				logger.Errorf("metric saving error: FreeMemory = %v", memoryStat.Free)
			}

			_, err = services.AddGaugeItem("CPUutilization1", model.Gauge(cpu.ClocksPerSec))
			if err != nil {
				logger.Errorf("metric saving error: CPUutilization1 = %v", cpu.ClocksPerSec)
			}
		}
	}
}

func collectMetricsNEW(ctx context.Context, wg *sync.WaitGroup, metricStat *runtime.MemStats, pollInterval int, services *service.MetricService, logger *zap.SugaredLogger) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("--> collect metrics goroutine exit")
			wg.Done()
			return

		case <-time.After(time.Duration(pollInterval) * time.Second):

			logger.Info("collecting default metrics")
			runtime.ReadMemStats(metricStat)

			for _, stat := range memStatsList {
				value, err := utils.GetFieldFromMemStats(metricStat, stat)
				if err != nil {
					logger.Errorf("%v: %s", err, stat)
				}
				_, err = services.AddGaugeItem(stat, value)
				if err != nil {
					logger.Errorf("metric saving error: %s = %v", stat, value)
				}
			}

			rndValue := rand.Intn(100)
			_, err := services.AddGaugeItem("RandomValue", model.Gauge(rndValue))
			if err != nil {
				logger.Errorf("metric saving error: RandomValue = %v", rndValue)
			}

			_, err = services.AddCounterItem("PollCount", model.Counter(1))
			if err != nil {
				logger.Errorf("metric saving error: PollCount")
			}
		}
	}
}

func pushMetricsToQueue(ctx context.Context, wg *sync.WaitGroup, cfg *agentConfig, services *service.MetricService, logger *zap.SugaredLogger) chan model.Metric {

	metricsQueue := make(chan model.Metric)

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("--> push metrics goroutine exit")
				close(metricsQueue)
				wg.Done()
				return

			case <-time.After(time.Duration(cfg.ReportInterval) * time.Second):
				metrics, err := services.GetAllMetrics()
				if err != nil {
					logger.Errorf("falied report: %v", err)
					continue
				}
				logger.Debug("push metrics to queue")
				for _, metric := range metrics {
					metricsQueue <- metric
				}
			}
		}
	}()

	return metricsQueue
}

func sendMetricsFromQueue(ctx context.Context, wg *sync.WaitGroup, worker int, queue <-chan model.Metric, cfg *agentConfig, services *service.MetricService, logger *zap.SugaredLogger) {
	for {
		select {
		case <-ctx.Done():
			defer wg.Done()
			logger.Infof("--> send metrics goroutine %d exit (%d)", worker, len(queue))
			if len(queue) > 0 {
				for metric := range queue {
					err := handlers.SendMetricFromQueue(services, &metric, cfg.Host, cfg.Key, logger)
					if err != nil {
						logger.Errorf("Falied to make request: %v", err)
					}
				}
			}
			return

		case metric := <-queue:
			err := handlers.SendMetricFromQueue(services, &metric, cfg.Host, cfg.Key, logger)
			if err != nil {
				logger.Errorf("Falied to make request: %v", err)
			}
		}
	}
}

func RunAgentAsync() error {

	cfgFlags := new(agentConfig)

	flag.StringVar(&cfgFlags.Host, "a", defServerHost, "server host")
	flag.StringVar(&cfgFlags.Key, "k", defAgentKey, "KEY")
	flag.IntVar(&cfgFlags.ReportInterval, "r", defReportInterval, "reportInterval")
	flag.IntVar(&cfgFlags.PollInterval, "p", defPollInterval, "pollInterval")
	flag.IntVar(&cfgFlags.RateLimit, "l", defRateLimit, "rateLimit")
	flag.Parse()

	cfg := initAgentConfigENV(cfgFlags)

	agentLogger, err := utils.InitLogger()
	if err != nil {
		log.Fatalf("agent logger init error: %v", err)
	}

	jsonConfig, _ := json.Marshal(cfg)
	agentLogger.Infof("Current agent config: %s", jsonConfig)

	db := storage.NewMemStorage()
	agentServices, err := service.NewMetricService(db, false, false, "")
	if err != nil {
		log.Fatalf("service construction error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	metricsStat := new(runtime.MemStats)

	var wg sync.WaitGroup

	wg.Add(1)
	go collectMetricsCPU(ctx, &wg, cfg.PollInterval, agentServices, agentLogger)

	wg.Add(1)
	go collectMetricsNEW(ctx, &wg, metricsStat, cfg.PollInterval, agentServices, agentLogger)

	wg.Add(1)
	metricsQueue := pushMetricsToQueue(ctx, &wg, cfg, agentServices, agentLogger)

	for worker := 1; worker <= cfg.RateLimit; worker++ {
		wg.Add(1)
		go sendMetricsFromQueue(ctx, &wg, worker, metricsQueue, cfg, agentServices, agentLogger)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	cancel()
	wg.Wait()

	agentLogger.Info("Agent shutdown gracefully")

	return nil
}
