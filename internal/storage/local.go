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

func (storage *MemStorage) Ping() error {
	return errors.New("its in-memory store")
}

func (storage *MemStorage) AddGaugeItem(key string, value model.Gauge) bool {
	storage.mx.Lock()
	defer storage.mx.Unlock()
	storage.GaugeItems[key] = value
	return true
}

func (storage *MemStorage) AddCounterItem(key string, value model.Counter) bool {
	storage.mx.Lock()
	defer storage.mx.Unlock()

	if _, ok := storage.CounterItems[key]; ok {
		storage.CounterItems[key] = storage.CounterItems[key] + value
	} else {
		storage.CounterItems[key] = value
	}
	return true
}

func (storage *MemStorage) AddMetricsPack(metrics *model.MetricsPack) bool {
	return false
}

func (storage *MemStorage) GetGaugeItem(key string) (model.Gauge, bool) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	if _, ok := storage.GaugeItems[key]; ok {
		return storage.GaugeItems[key], true
	}
	return 0, false
}

func (storage *MemStorage) GetCounterItem(key string) (model.Counter, bool) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	if _, ok := storage.CounterItems[key]; ok {
		return storage.CounterItems[key], true
	}
	return 0, false
}

func (storage *MemStorage) ResetCounterItem(key string) bool {
	storage.mx.RLock()
	defer storage.mx.RUnlock()

	if _, ok := storage.CounterItems[key]; ok {
		storage.CounterItems[key] = model.Counter(0)
		return true
	}
	return false
}

func (storage *MemStorage) GetGaugeItems() (map[string]model.Gauge, bool) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()
	result := storage.GaugeItems
	return result, true
}

func (storage *MemStorage) GetCounterItems() (map[string]model.Counter, bool) {
	storage.mx.RLock()
	defer storage.mx.RUnlock()
	result := storage.CounterItems
	return result, true
}
