package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

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

	GetGaugeItems() (map[string]model.Gauge, bool)
	GetCounterItems() (map[string]model.Counter, bool)

	GetStringGaugeItems() (map[string]string, bool)
	GetStringCounterItems() (map[string]string, bool)
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

func (h *MetricService) GetGaugeItems() (map[string]model.Gauge, error) {
	if result, ok := h.store.GetGaugeItems(); ok {
		return result, nil
	}
	return map[string]model.Gauge{}, ErrorGettingMetrics
}

func (h *MetricService) GetCounterItems() (map[string]model.Counter, error) {
	if result, ok := h.store.GetCounterItems(); ok {
		return result, nil
	}
	return map[string]model.Counter{}, ErrorGettingMetrics
}

func (h *MetricService) GetStringGaugeItems() (map[string]string, error) {
	result, err := h.store.GetStringGaugeItems()
	if !err {
		return map[string]string{}, ErrorGettingMetrics
	}
	return result, nil
}

func (h *MetricService) GetStringCounterItems() (map[string]string, error) {
	result, err := h.store.GetStringCounterItems()
	if !err {
		return map[string]string{}, ErrorGettingMetrics
	}
	return result, nil
}

type metricsBackup struct {
	Metrics []model.Metric `json:"metrics"`
}

func (h *MetricService) ExportToJSON() ([]byte, error) {

	var metricOut metricsBackup

	counter, err := h.GetCounterItems()
	if err != nil {
		return nil, err
	}

	for key, value := range counter {
		metricValue := int64(value)
		metric := model.Metric{
			ID:    key,
			MType: "counter",
			Delta: &metricValue,
		}

		metricOut.Metrics = append(metricOut.Metrics, metric)

	}

	gauge, err := h.GetGaugeItems()
	if err != nil {
		return nil, err
	}

	for key, value := range gauge {

		metricValue := float64(value)
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
	var metricStruct metricsBackup

	err := json.Unmarshal(data, &metricStruct)
	if err != nil {
		return err
	}

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
