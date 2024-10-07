package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bbquite/mca-server/internal/model"
	"github.com/bbquite/mca-server/internal/utils"
	"go.uber.org/zap"
)

type MemStorageRepo interface {
	AddGaugeItem(key string, value model.Gauge) error
	AddCounterItem(key string, value model.Counter) error

	AddMetricsPack(metrics *model.MetricsPack) error

	GetGaugeItem(key string) (model.Gauge, error)
	GetCounterItem(key string) (model.Counter, error)
	ResetCounterItem(key string) error

	GetGaugeItems() (map[string]model.Gauge, error)
	GetCounterItems() (map[string]model.Counter, error)

	Ping() error
}

type MetricService struct {
	store           MemStorageRepo
	syncSave        bool
	filePath        string
	isDatabaseUsage bool
	logger          *zap.SugaredLogger
}

func NewMetricService(store MemStorageRepo, syncSave bool, isDatabaseUsage bool, filePath string) (*MetricService, error) {

	logger, err := utils.InitLogger()
	if err != nil {
		return nil, err
	}

	return &MetricService{
		store:           store,
		syncSave:        syncSave,
		filePath:        filePath,
		isDatabaseUsage: isDatabaseUsage,
		logger:          logger,
	}, nil
}

func (s *MetricService) PingDatabase() error {
	err := s.store.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *MetricService) AddGaugeItem(key string, value model.Gauge) (model.Gauge, error) {
	err := s.store.AddGaugeItem(key, value)
	if err != nil {
		return 0, err
	}

	if s.syncSave {
		err = s.SaveToFile(s.filePath)
		if err != nil {
			s.logger.Error(err)
		}
	}
	return model.Gauge(value), nil
}

func (s *MetricService) AddCounterItem(key string, value model.Counter) (model.Counter, error) {
	err := s.store.AddCounterItem(key, value)
	if err != nil {
		return 0, err
	}

	if s.syncSave {
		err = s.SaveToFile(s.filePath)
		if err != nil {
			s.logger.Error(err)
		}
	}
	return model.Counter(value), nil
}

func (s *MetricService) GetGaugeItem(key string) (model.Gauge, error) {
	item, err := s.store.GetGaugeItem(key)
	if err != nil {
		return 0, err
	}

	return item, nil
}

func (s *MetricService) GetCounterItem(key string) (model.Counter, error) {
	item, err := s.store.GetCounterItem(key)
	if err != nil {
		return 0, err
	}

	return item, nil
}

func (s *MetricService) ResetCounterItem(key string) error {
	err := s.store.ResetCounterItem(key)
	if err != nil {
		return err
	}
	return nil
}

func (s *MetricService) GetGaugeItems() (map[string]model.Gauge, error) {
	items, err := s.store.GetGaugeItems()
	if err != nil {
		return map[string]model.Gauge{}, err
	}

	return items, nil
}

func (s *MetricService) GetCounterItems() (map[string]model.Counter, error) {

	items, err := s.store.GetCounterItems()
	if err != nil {
		return map[string]model.Counter{}, err
	}

	return items, nil
}

func (s *MetricService) GetAllMetrics() (model.MetricsPack, error) {
	var metricResult model.MetricsPack

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

		metricResult = append(metricResult, metric)

	}

	gauge, err := s.GetGaugeItems()
	if err != nil {
		return metricResult, err
	}

	for key, value := range gauge {

		metricValue := float64(value)
		metric := model.Metric{
			ID:    key,
			MType: "gauge",
			Value: &metricValue,
		}

		metricResult = append(metricResult, metric)
	}

	return metricResult, nil
}

func (s *MetricService) ExportToJSON() ([]byte, error) {

	metricsPack, err := s.GetAllMetrics()
	if err != nil {
		return nil, err
	}

	metricsJSON, err := json.Marshal(metricsPack)
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

	for _, element := range metricStruct {
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
	s.logger.Infof("save data: %s", data)
	if err != nil {
		return err
	}
	os.WriteFile(filePath, data, 0666)

	return nil
}

func (s *MetricService) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	s.logger.Infof("load data: %s", data)
	if err != nil {
		return err
	}
	if err := s.ImportFromJSON(data); err != nil {
		return err
	}
	return nil
}
