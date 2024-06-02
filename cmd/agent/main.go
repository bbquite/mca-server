package main

import (
	"flag"
	"github.com/bbquite/mca-server/internal/handlers"
	"github.com/bbquite/mca-server/internal/service"
	"github.com/bbquite/mca-server/internal/storage"
	"github.com/joho/godotenv"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
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

func agentRun(opt *Options) error {

	log.Printf("INFO | Server HOST: %s", opt.a)

	db := storage.NewMemStorage()
	agentServices := service.NewMetricService(db)
	memStat := new(runtime.MemStats)

	collectTicker := time.NewTicker(time.Duration(opt.p) * time.Second)
	requestTicker := time.NewTicker(time.Duration(opt.p) * time.Second)

	for {
		select {
		case <-collectTicker.C:
			handlers.MetricsCollect(memStat, agentServices)
		case <-requestTicker.C:
			//err := handlers.MetricsUriRequest(agentServices, opt.a)
			err := handlers.MetricsPostRequest(agentServices, opt.a)
			if err != nil {
				log.Printf("ERROR | Falied to make request: %v", err)
			}
		}
	}
}

func main() {
	opt := initOptions()

	if err := agentRun(opt); err != nil {
		log.Fatal(err)
	}
}
