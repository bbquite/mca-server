package service

import (
	"errors"
	"github.com/bbquite/mca-server/internal/model"
)

var (
	ErrorAddingGauge   = errors.New("no gauge value added")
	ErrorAddingCounter = errors.New("no counter value added")

	ErrorGaugeNotFound   = errors.New("gauge not found")
	ErrorCounterNotFound = errors.New("counters not found")
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
		return map[string]string{}, errors.New("error getting metrics")

	}
	return result, nil
}

func (h *MetricService) GetAllCounterItems() (map[string]string, error) {
	result, err := h.store.GetAllCounterItems()
	if !err {
		return map[string]string{}, errors.New("error getting metrics")
	}
	return result, nil
}
