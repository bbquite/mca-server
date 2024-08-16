package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bbquite/mca-server/internal/model"
)

var (
	ErrorAddingGauge   = errors.New("no gauge value added")
	ErrorAddingCounter = errors.New("no counter value added")

	ErrorGaugeNotFound   = errors.New("gauge not found")
	ErrorCounterNotFound = errors.New("counters not found")
	ErrorGettingMetrics  = errors.New("error getting metrics")
)

type MemStorageRepo interface {
	AddGaugeItem(key string, value model.Gauge) bool
	AddCounterItem(key string, value model.Counter) bool

	GetGaugeItem(key string) (model.Gauge, bool)
	GetCounterItem(key string) (model.Counter, bool)

	GetAllGaugeItems() (map[string]string, bool)
	GetAllCounterItems() (map[string]string, bool)
}

type MetricService struct {
	store MemStorageRepo
}

func NewMetricService(store MemStorageRepo) *MetricService {
	return &MetricService{store: store}
}

func (h *MetricService) AddGaugeItem(key string, value model.Gauge) (model.Gauge, error) {
	if ok := h.store.AddGaugeItem(key, value); ok {
		return model.Gauge(value), nil
	}
	return 0, ErrorAddingGauge
}

func (h *MetricService) AddCounterItem(key string, value model.Counter) (model.Counter, error) {
	if ok := h.store.AddCounterItem(key, value); ok {
		return model.Counter(value), nil
	}
	return 0, ErrorAddingCounter
}

func (h *MetricService) GetGaugeItem(key string) (model.Gauge, error) {
	if gauge, ok := h.store.GetGaugeItem(key); ok {
		return gauge, nil
	}
	return 0, ErrorGaugeNotFound
}

func (h *MetricService) GetCounterItem(key string) (model.Counter, error) {
	if counter, ok := h.store.GetCounterItem(key); ok {
		return counter, nil
	}
	return 0, ErrorCounterNotFound
}

func (h *MetricService) GetAllGaugeItems() (map[string]string, error) {
	result, err := h.store.GetAllGaugeItems()
	if !err {
		return map[string]string{}, ErrorGettingMetrics
	}
	return result, nil
}

func (h *MetricService) GetAllCounterItems() (map[string]string, error) {
	result, err := h.store.GetAllCounterItems()
	if !err {
		return map[string]string{}, ErrorGettingMetrics
	}
	return result, nil
}

type metricsOutput struct {
	Metrics []model.Metric `json:"metrics"`
}

func (h *MetricService) ExportToJSON() ([]byte, error) {

	var metricOut metricsOutput

	counter, err := h.GetAllCounterItems()
	if err != nil {
		return nil, err
	}

	for key, value := range counter {
		metricValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse int err: %v", err)
		}

		metric := model.Metric{
			ID:    key,
			MType: "counter",
			Delta: &metricValue,
		}

		metricOut.Metrics = append(metricOut.Metrics, metric)

	}

	gauge, err := h.GetAllGaugeItems()
	if err != nil {
		return nil, err
	}

	for key, value := range gauge {

		metricValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("parse float err: %v", err)
		}

		metric := model.Metric{
			ID:    key,
			MType: "gauge",
			Value: &metricValue,
		}

		metricOut.Metrics = append(metricOut.Metrics, metric)
	}

	metricsJSON, err := json.Marshal(metricOut)
	if err != nil {
		return nil, fmt.Errorf("err: %v", err)
	}

	return metricsJSON, nil
}

func (h *MetricService) ImportFromJSON(data []byte) error {
	log.Printf("%s", data)
	var metricStruct metricsOutput

	err := json.Unmarshal(data, &metricStruct)
	if err != nil {
		return err
	}

	log.Print(metricStruct)

	for _, element := range metricStruct.Metrics {
		switch element.MType {
		case "gauge":
			_, err = h.AddGaugeItem(element.ID, model.Gauge(*element.Value))
			if err != nil {
				return err
			}

		case "counter":
			_, err = h.AddCounterItem(element.ID, model.Counter(*element.Delta))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *MetricService) SaveToFile(filePath string) error {
	data, err := h.ExportToJSON()
	log.Printf("save data: %s", data)
	if err != nil {
		return err
	}
	os.WriteFile(filePath, data, 0666)

	return nil
}

func (h *MetricService) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	log.Printf("load data: %s", data)
	if err != nil {
		return err
	}
	if err := h.ImportFromJSON(data); err != nil {
		return err
	}
	return nil
}
