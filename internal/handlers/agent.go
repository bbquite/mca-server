package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/service"
	"go.uber.org/zap"
)

func MetricsCollect(memStat *runtime.MemStats, services *service.MetricService, logger *zap.SugaredLogger) {

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

func MetricsURIRequest(services *service.MetricService, host string) error {
	var url string

	log.Printf("INFO | Sending metrics to %s", host)

	client := http.Client{
		Timeout: time.Second * 1,
	}

	gauge, err := services.GetAllGaugeItems()
	if err != nil {
		return fmt.Errorf("gauge metrics collection error: %v", err)
	}
	for key, value := range gauge {
		url = fmt.Sprintf("http://%s/update/%s/%s/%s", host, "gauge", key, value)

		request, _ := http.NewRequest(http.MethodPost, url, http.NoBody)
		request.Header.Set("Content-Type", "Content-Type: text/plain")

		response, err := client.Do(request)
		if err != nil {
			log.Printf("clent request error: %v", err)
			return nil
		}

		err = response.Body.Close()
		if err != nil {
			return err
		}

		if response.StatusCode != 200 {
			return fmt.Errorf("bad server response, status code: %d", response.StatusCode)
		}
	}

	counter, err := services.GetAllCounterItems()
	if err != nil {
		return fmt.Errorf("counter metrics collection error: %v", err)
	}
	for key, value := range counter {
		url = fmt.Sprintf("http://%s/update/%s/%s/%s", host, "counter", key, value)

		request, _ := http.NewRequest(http.MethodPost, url, http.NoBody)
		request.Header.Set("Content-Type", "Content-Type: text/plain")

		response, err := client.Do(request)
		if err != nil {
			log.Printf("clent request error: %v", err)
			return nil
		}

		err = response.Body.Close()
		if err != nil {
			return err
		}

		if response.StatusCode != 200 {
			return fmt.Errorf("bad server response, status code: %d", response.StatusCode)
		}
	}

	return nil
}

func MetricsJSONRequest(services *service.MetricService, host string, logger *zap.SugaredLogger) error {

	url := fmt.Sprintf("http://%s/update/", host)
	client := http.Client{}

	counter, err := services.GetAllCounterItems()
	logger.Infof("sending counter metrics")
	if err != nil {
		return fmt.Errorf("counter metrics sending error: %v", err)
	}

	for key, value := range counter {
		metricValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int err: %v", err)
		}

		metric := model.Metric{
			ID:    key,
			MType: "counter",
			Delta: &metricValue,
		}

		body, err := json.Marshal(metric)
		if err != nil {
			return fmt.Errorf("err: %v", err)
		}

		bodySend := bytes.NewBuffer(body)

		request, err := http.NewRequest(http.MethodPost, url, bodySend)
		if err != nil {
			logger.Error(err)
			return err
		}

		request.Header.Set("Content-Type", "application/json")

		logger.Debugf("TRY %s %s", url, request.Body)

		response, err := client.Do(request)
		if err != nil {
			logger.Error(err)
			return nil
		}

		defer response.Body.Close()
		logger.Debugf("OK %s %s", url, response.Status)

		if response.StatusCode != 200 {
			return fmt.Errorf("bad server response, status code: %d", response.StatusCode)
		}
	}

	gauge, err := services.GetAllGaugeItems()
	logger.Infof("sending gauge metrics")
	if err != nil {
		return fmt.Errorf("gauge metrics sending error: %v", err)
	}

	for key, value := range gauge {

		metricValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("parse float err: %v", err)
		}

		metric := model.Metric{
			ID:    key,
			MType: "gauge",
			Value: &metricValue,
		}

		body, err := json.Marshal(metric)
		if err != nil {
			logger.Error(err)
			return nil
		}

		bodySend := bytes.NewBuffer(body)

		request, err := http.NewRequest(http.MethodPost, url, bodySend)
		if err != nil {
			logger.Error(err)
			return err
		}

		request.Header.Set("Content-Type", "application/json")

		logger.Debugf("TRY %s %s", url, request.Body)

		response, err := client.Do(request)
		if err != nil {
			logger.Error(err)
			return nil
		}

		defer response.Body.Close()
		logger.Debugf("OK %s %s", url, response.Status)

		if response.StatusCode != 200 {
			return fmt.Errorf("bad server response, status code: %d", response.StatusCode)
		}
	}

	return nil
}
