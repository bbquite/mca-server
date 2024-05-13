package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func CollectMemStat(memStat runtime.MemStats) map[string]string {
	result := map[string]string{
		"Alloc":         strconv.FormatUint(memStat.Alloc, 10),
		"BuckHashSys":   strconv.FormatUint(memStat.BuckHashSys, 10),
		"Frees":         strconv.FormatUint(memStat.Frees, 10),
		"GCCPUFraction": strconv.FormatFloat(memStat.GCCPUFraction, 'f', 2, 64),
		"GCSys":         strconv.FormatUint(memStat.GCSys, 10),
		"HeapAlloc":     strconv.FormatUint(memStat.HeapAlloc, 10),
		"HeapIdle":      strconv.FormatUint(memStat.HeapIdle, 10),
		"HeapInuse":     strconv.FormatUint(memStat.HeapInuse, 10),
		"HeapObjects":   strconv.FormatUint(memStat.HeapObjects, 10),
		"HeapReleased":  strconv.FormatUint(memStat.HeapReleased, 10),
		"HeapSys":       strconv.FormatUint(memStat.HeapSys, 10),
		"LastGC":        strconv.FormatUint(memStat.LastGC, 10),
		"Lookups":       strconv.FormatUint(memStat.Lookups, 10),
		"MCacheInuse":   strconv.FormatUint(memStat.MCacheInuse, 10),
		"MCacheSys":     strconv.FormatUint(memStat.MCacheSys, 10),
		"MSpanInuse":    strconv.FormatUint(memStat.MSpanInuse, 10),
		"Mallocs":       strconv.FormatUint(memStat.Mallocs, 10),
		"NextGC":        strconv.FormatUint(memStat.NextGC, 10),
		"NumForcedGC":   strconv.FormatUint(uint64(memStat.NumForcedGC), 10),
		"NumGC":         strconv.FormatUint(uint64(memStat.NumGC), 10),
		"OtherSys":      strconv.FormatUint(memStat.OtherSys, 10),
		"PauseTotalNs":  strconv.FormatUint(memStat.PauseTotalNs, 10),
		"StackInuse":    strconv.FormatUint(memStat.StackInuse, 10),
		"StackSys":      strconv.FormatUint(memStat.StackSys, 10),
		"Sys":           strconv.FormatUint(memStat.Sys, 10),
		"TotalAlloc":    strconv.FormatUint(memStat.TotalAlloc, 10),
		"RandomValue":   strconv.Itoa(rand.Intn(100)),
	}
	return result
}

func SendRequestMemStat(host string, mType string, mName string, mValue string) error {

	url := fmt.Sprintf("http://%s/update/%s/%s/%s", host, mType, mName, mValue)

	client := http.Client{
		Timeout: time.Second * 1,
	}

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

	log.Printf("send %s, status code %d", url, response.StatusCode)
	return nil
}
