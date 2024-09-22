package storage

import (
	"errors"
	"sync"

	"github.com/bbquite/mca-server/internal/model"
)

type MemStorage struct {
	GaugeItems   map[string]model.Gauge
	CounterItems map[string]model.Counter
	mx           sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeItems:   make(map[string]model.Gauge),
		CounterItems: make(map[string]model.Counter),
	}
}

func (storage *MemStorage) AddGaugeItem(key string, value model.Gauge) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.GaugeItems[key] = value
	return nil
}

func (storage *MemStorage) AddCounterItem(key string, value model.Counter) error {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	if _, ok := storage.CounterItems[key]; ok {
		storage.CounterItems[key] = storage.CounterItems[key] + value
	} else {
		storage.CounterItems[key] = value
	}
	return nil
}

func (storage *MemStorage) GetGaugeItem(key string) (model.Gauge, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	if _, ok := storage.GaugeItems[key]; ok {
		return storage.GaugeItems[key], nil
	}
	return 0, ErrorGaugeNotFound
}

func (storage *MemStorage) GetCounterItem(key string) (model.Counter, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	if _, ok := storage.CounterItems[key]; ok {
		return storage.CounterItems[key], nil
	}
	return 0, ErrorCounterNotFound
}

func (storage *MemStorage) GetGaugeItems() (map[string]model.Gauge, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()
	result := storage.GaugeItems
	return result, nil
}

func (storage *MemStorage) GetCounterItems() (map[string]model.Counter, error) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()
	result := storage.CounterItems
	return result, nil
}

func (storage *MemStorage) ResetCounterItem(key string) error {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	if _, ok := storage.CounterItems[key]; ok {
		storage.CounterItems[key] = model.Counter(0)
		return ErrorResetCounter
	}
	return nil
}

func (storage *MemStorage) Ping() error {
	return errors.New("UNUSED")
}

func (storage *MemStorage) AddMetricsPack(metrics *model.MetricsPack) error {
	return errors.New("UNUSED")
}
