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

func collectRun() {
	for {
		log.Print(1)
		time.Sleep(2 * time.Second)
	}

}

func collectRun2() {
	for {
		log.Print(2)
		time.Sleep(5 * time.Second)
	}

}

func agentRun(opt *Options) error {

	log.Printf("Server HOST: %s", opt.a)

	db := storage.NewMemStorage()
	agentServices := service.NewMetricService(db)
	memStat := new(runtime.MemStats)

	go handlers.CollectStat(memStat, agentServices, opt.p)

	err := handlers.MakeRequestStat(agentServices, opt.a, opt.r)
	if err != nil {
		log.Print(err)
	}

	return nil
}

func main() {
	opt := initOptions()

	if err := agentRun(opt); err != nil {
		log.Fatal(err)
	}
}
