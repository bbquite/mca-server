package storage

import (
	"github.com/bbquite/mca-server/internal/model"
)

type MemStorage struct {
	GaugeItems   map[string]model.Gauge
	CounterItems map[string]model.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeItems:   make(map[string]model.Gauge),
		CounterItems: make(map[string]model.Counter),
	}
}

func (storage *MemStorage) AddGaugeItem(key string, value model.Gauge) bool {
	storage.GaugeItems[key] = value
	return true
}

func (storage *MemStorage) AddCounterItem(key string, value model.Counter) bool {
	storage.CounterItems[key] = value
	return true
}

func (storage *MemStorage) GetGaugeItem(key string, value model.Gauge) (model.Gauge, bool) {
	if _, ok := storage.GaugeItems[key]; ok {
		return storage.GaugeItems[key], true
	}
	return 0, false
}

func (storage *MemStorage) GetCounterItem(key string, value model.Counter) (model.Counter, bool) {
	if _, ok := storage.CounterItems[key]; ok {
		return storage.CounterItems[key], true
	}
	return 0, false
}

func (storage *MemStorage) GetAllGaugeItems() (map[string]model.Gauge, bool) {
	result := storage.GaugeItems
	return result, true
}

func (storage *MemStorage) GetAllCounterItems() (map[string]model.Counter, bool) {
	result := storage.CounterItems
	return result, true
}
