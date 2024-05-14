package handlers

import (
	"fmt"
	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/service"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func CollectStat(memStat *runtime.MemStats, services *service.MetricService, sleepDuration int) {

	for {
		time.Sleep(time.Duration(sleepDuration) * time.Second)
		log.Print("Collecting metrics and recording in storage")
		runtime.ReadMemStats(memStat)

		_, err := services.AddGaugeItem("Alloc", model.Gauge(memStat.Alloc))
		if err != nil {
			log.Printf("metric saving error: Alloc = %v", memStat.Alloc)
		}

		_, err = services.AddGaugeItem("BuckHashSys", model.Gauge(memStat.BuckHashSys))
		if err != nil {
			log.Printf("metric saving error: BuckHashSys = %v", memStat.BuckHashSys)
		}

		_, err = services.AddGaugeItem("Frees", model.Gauge(memStat.Frees))
		if err != nil {
			log.Printf("metric saving error: Frees = %v", memStat.Frees)
		}

		_, err = services.AddGaugeItem("GCCPUFraction", model.Gauge(memStat.GCCPUFraction))
		if err != nil {
			log.Printf("metric saving error: Alloc = %v", memStat.GCCPUFraction)
		}

		_, err = services.AddGaugeItem("GCSys", model.Gauge(memStat.GCSys))
		if err != nil {
			log.Printf("metric saving error: GCSys = %v", memStat.GCSys)
		}

		_, err = services.AddGaugeItem("HeapAlloc", model.Gauge(memStat.HeapAlloc))
		if err != nil {
			log.Printf("metric saving error: HeapAlloc = %v", memStat.HeapAlloc)
		}

		_, err = services.AddGaugeItem("HeapIdle", model.Gauge(memStat.HeapIdle))
		if err != nil {
			log.Printf("metric saving error: HeapIdle = %v", memStat.HeapIdle)
		}

		_, err = services.AddGaugeItem("HeapInuse", model.Gauge(memStat.HeapInuse))
		if err != nil {
			log.Printf("metric saving error: HeapInuse = %v", memStat.HeapInuse)
		}

		_, err = services.AddGaugeItem("HeapObjects", model.Gauge(memStat.HeapObjects))
		if err != nil {
			log.Printf("metric saving error: HeapObjects = %v", memStat.HeapObjects)
		}

		_, err = services.AddGaugeItem("HeapReleased", model.Gauge(memStat.HeapReleased))
		if err != nil {
			log.Printf("metric saving error: HeapReleased = %v", memStat.HeapReleased)
		}

		_, err = services.AddGaugeItem("HeapSys", model.Gauge(memStat.HeapSys))
		if err != nil {
			log.Printf("metric saving error: HeapSys = %v", memStat.HeapSys)
		}

		_, err = services.AddGaugeItem("LastGC", model.Gauge(memStat.LastGC))
		if err != nil {
			log.Printf("metric saving error: LastGC = %v", memStat.LastGC)
		}

		_, err = services.AddGaugeItem("Lookups", model.Gauge(memStat.Lookups))
		if err != nil {
			log.Printf("metric saving error: Lookups = %v", memStat.Lookups)
		}

		_, err = services.AddGaugeItem("MCacheInuse", model.Gauge(memStat.MCacheInuse))
		if err != nil {
			log.Printf("metric saving error: MCacheInuse = %v", memStat.MCacheInuse)
		}

		_, err = services.AddGaugeItem("MCacheSys", model.Gauge(memStat.MCacheSys))
		if err != nil {
			log.Printf("metric saving error: MSpanInuse = %v", memStat.MSpanInuse)
		}

		_, err = services.AddGaugeItem("Mallocs", model.Gauge(memStat.Mallocs))
		if err != nil {
			log.Printf("metric saving error: Mallocs = %v", memStat.Mallocs)
		}

		_, err = services.AddGaugeItem("NextGC", model.Gauge(memStat.NextGC))
		if err != nil {
			log.Printf("metric saving error: NextGC = %v", memStat.NextGC)
		}

		_, err = services.AddGaugeItem("NumForcedGC", model.Gauge(memStat.NumForcedGC))
		if err != nil {
			log.Printf("metric saving error: NumForcedGC = %v", memStat.NumForcedGC)
		}

		_, err = services.AddGaugeItem("NumGC", model.Gauge(memStat.NumGC))
		if err != nil {
			log.Printf("metric saving error: NumGC = %v", memStat.NumGC)
		}

		_, err = services.AddGaugeItem("OtherSys", model.Gauge(memStat.OtherSys))
		if err != nil {
			log.Printf("metric saving error: OtherSys = %v", memStat.OtherSys)
		}

		_, err = services.AddGaugeItem("PauseTotalNs", model.Gauge(memStat.PauseTotalNs))
		if err != nil {
			log.Printf("metric saving error: PauseTotalNs = %v", memStat.PauseTotalNs)
		}

		_, err = services.AddGaugeItem("StackInuse", model.Gauge(memStat.StackInuse))
		if err != nil {
			log.Printf("metric saving error: StackInuse = %v", memStat.StackInuse)
		}

		_, err = services.AddGaugeItem("StackSys", model.Gauge(memStat.StackSys))
		if err != nil {
			log.Printf("metric saving error: StackSys = %v", memStat.StackSys)
		}

		_, err = services.AddGaugeItem("Sys", model.Gauge(memStat.Sys))
		if err != nil {
			log.Printf("metric saving error: NextGC = %v", memStat.Sys)
		}

		_, err = services.AddGaugeItem("TotalAlloc", model.Gauge(memStat.TotalAlloc))
		if err != nil {
			log.Printf("metric saving error: NextGC = %v", memStat.TotalAlloc)
		}

		rndValue := rand.Intn(100)
		_, err = services.AddGaugeItem("RandomValue", model.Gauge(rndValue))
		if err != nil {
			log.Printf("metric saving error: RandomValue = %v", rndValue)
		}

		_, err = services.AddCounterItem("PollCount", model.Counter(1))
		if err != nil {
			log.Printf("metric saving error: PollCount")
		}
	}
}

func MakeRequestStat(services *service.MetricService, host string, sleepDuration int) error {
	var url string
	for {

		time.Sleep(time.Duration(sleepDuration) * time.Second)
		log.Printf("Sending metrics to %s", host)

		client := http.Client{
			Timeout: time.Second * 1,
		}

		gauge, _ := services.GetAllGaugeItems()
		for key, value := range gauge {
			url = fmt.Sprintf("http://%s/update/%s/%s/%s", host, "gauge", key, value)

			request, _ := http.NewRequest(http.MethodPost, url, http.NoBody)
			request.Header.Set("Content-Type", "Content-Type: text/plain")

			response, err := client.Do(request)
			if err != nil {
				return err
			}
			err = response.Body.Close()
			if err != nil {
				return err
			}
		}

		counter, _ := services.GetAllCounterItems()
		for key, value := range counter {
			url = fmt.Sprintf("http://%s/update/%s/%s/%s", host, "counter", key, value)

			request, _ := http.NewRequest(http.MethodPost, url, http.NoBody)
			request.Header.Set("Content-Type", "Content-Type: text/plain")

			response, err := client.Do(request)
			if err != nil {
				return err
			}
			err = response.Body.Close()
			if err != nil {
				return err
			}
		}
	}
}
