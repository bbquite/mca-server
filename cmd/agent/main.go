package main

import (
	"flag"
	"github.com/bbquite/mca-server/internal/handlers"
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

var (
	PollCount = 0
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
	memStat := new(runtime.MemStats)

	log.Printf("Server HOST: %s", opt.a)

	for {
		runtime.ReadMemStats(memStat)
		memStatMap := handlers.CollectMemStat(*memStat)
		for key, value := range memStatMap {
			err := handlers.SendRequestMemStat(opt.a, "gauge", key, value)
			if err != nil {
				log.Print(err)
				continue
			}
		}

		PollCount += 1

		err := handlers.SendRequestMemStat(opt.a, "counter", "PollCount", strconv.Itoa(PollCount))
		if err != nil {
			log.Print(err)
		}

		time.Sleep(2 * time.Second)
	}
}

func main() {
	opt := initOptions()

	if err := agentRun(opt); err != nil {
		log.Fatal(err)
	}
}
