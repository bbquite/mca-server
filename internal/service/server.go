package service

import (
	"errors"
	"github.com/bbquite/mca-server/internal/model"
)

var (
	ErrorAddingGauge   = errors.New("no gauge value added")
	ErrorAddingCounter = errors.New("no counter value added")

	ErrorGaugeNotFound    = errors.New("gauge not found")
	ErrorCountersNotFound = errors.New("counters not found")
)

type MemStorageRepo interface {
	AddGaugeItem(key string, value model.Gauge) bool
	AddCounterItem(key string, value model.Counter) bool

	GetGaugeItem(key string, value model.Gauge) (model.Gauge, bool)
	GetCounterItem(key string, value model.Counter) (model.Counter, bool)

	GetAllGaugeItems() (map[string]model.Gauge, bool)
	GetAllCounterItems() (map[string]model.Counter, bool)
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

func (h *MetricService) GetGaugeItem(key string, value model.Gauge) (model.Gauge, error) {
	if gauge, ok := h.store.GetGaugeItem(key, value); ok {
		return gauge, nil
	}
	return 0, ErrorGaugeNotFound
}

func (h *MetricService) GetCounterItem(key string, value model.Counter) (model.Counter, error) {
	if counter, ok := h.store.GetCounterItem(key, value); ok {
		return counter, nil
	}
	return 0, ErrorCountersNotFound
}

func (h *MetricService) GetAllGaugeItems() (map[string]model.Gauge, error) {
	if result, ok := h.store.GetAllGaugeItems(); ok {
		return result, nil
	}

	return make(map[string]model.Gauge), errors.New("test error")
}

func (h *MetricService) GetAllCounterItems() (map[string]model.Counter, error) {
	if result, ok := h.store.GetAllCounterItems(); ok {
		return result, nil
	}

	return make(map[string]model.Counter), errors.New("test error")
}
