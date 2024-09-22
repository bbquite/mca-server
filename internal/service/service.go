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

	AddMetricsPack(metrics *model.MetricsPack) bool

	GetGaugeItem(key string) (model.Gauge, bool)
	GetCounterItem(key string) (model.Counter, bool)
	ResetCounterItem(key string) bool

	GetGaugeItems() (map[string]model.Gauge, bool)
	GetCounterItems() (map[string]model.Counter, bool)

	GetStringGaugeItems() (map[string]string, bool)
	GetStringCounterItems() (map[string]string, bool)

	Ping() error
}

type MetricService struct {
	store           MemStorageRepo
	syncSave        bool
	filePath        string
	isDatabaseUsage bool
}

func NewMetricService(store MemStorageRepo, syncSave bool, isDatabaseUsage bool, filePath string) *MetricService {
	return &MetricService{
		store:           store,
		syncSave:        syncSave,
		filePath:        filePath,
		isDatabaseUsage: isDatabaseUsage,
	}
}

func (s *MetricService) PingDatabase() error {
	err := s.store.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *MetricService) AddGaugeItem(key string, value model.Gauge) (model.Gauge, error) {
	if ok := s.store.AddGaugeItem(key, value); ok {
		if s.syncSave {
			s.SaveToFile(s.filePath)
		}
		return model.Gauge(value), nil
	}
	return 0, ErrorAddingGauge
}

func (s *MetricService) AddCounterItem(key string, value model.Counter) (model.Counter, error) {
	if ok := s.store.AddCounterItem(key, value); ok {
		if s.syncSave {
			s.SaveToFile(s.filePath)
		}
		return model.Counter(value), nil
	}
	return 0, ErrorAddingCounter
}

func (s *MetricService) GetGaugeItem(key string) (model.Gauge, error) {
	if gauge, ok := s.store.GetGaugeItem(key); ok {
		return gauge, nil
	}
	return 0, ErrorGaugeNotFound
}

func (s *MetricService) GetCounterItem(key string) (model.Counter, error) {
	if counter, ok := s.store.GetCounterItem(key); ok {
		return counter, nil
	}
	return 0, ErrorCounterNotFound
}

func (s *MetricService) ResetCounterItem(key string) error {
	if ok := s.store.ResetCounterItem(key); ok {
		return nil
	}
	return ErrorCounterNotFound
}

func (s *MetricService) GetGaugeItems() (map[string]model.Gauge, error) {
	if result, ok := s.store.GetGaugeItems(); ok {
		return result, nil
	}
	return map[string]model.Gauge{}, ErrorGettingMetrics
}

func (s *MetricService) GetCounterItems() (map[string]model.Counter, error) {
	if result, ok := s.store.GetCounterItems(); ok {
		return result, nil
	}
	return map[string]model.Counter{}, ErrorGettingMetrics
}

func (s *MetricService) GetStringGaugeItems() (map[string]string, error) {
	result, err := s.store.GetStringGaugeItems()
	if !err {
		return map[string]string{}, ErrorGettingMetrics
	}
	return result, nil
}

func (s *MetricService) GetStringCounterItems() (map[string]string, error) {
	result, err := s.store.GetStringCounterItems()
	if !err {
		return map[string]string{}, ErrorGettingMetrics
	}
	return result, nil
}

func (s *MetricService) ExportToJSON() ([]byte, error) {

	var metricOut model.MetricsPack

	counter, err := s.GetCounterItems()
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

	gauge, err := s.GetGaugeItems()
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

func (s *MetricService) ImportFromJSON(data []byte) error {
	var metricStruct model.MetricsPack

	err := json.Unmarshal(data, &metricStruct)
	if err != nil {
		return err
	}

	if s.isDatabaseUsage {
		s.store.AddMetricsPack(&metricStruct)
		return nil
	}

	for _, element := range metricStruct.Metrics {
		switch element.MType {
		case "gauge":
			_, err = s.AddGaugeItem(element.ID, model.Gauge(*element.Value))
			if err != nil {
				return err
			}

		case "counter":
			_, err = s.AddCounterItem(element.ID, model.Counter(*element.Delta))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *MetricService) SaveToFile(filePath string) error {
	data, err := s.ExportToJSON()
	log.Printf("save data: %s", data)
	if err != nil {
		return err
	}
	os.WriteFile(filePath, data, 0666)

	return nil
}

func (s *MetricService) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	log.Printf("load data: %s", data)
	if err != nil {
		return err
	}
	if err := s.ImportFromJSON(data); err != nil {
		return err
	}
	return nil
}
