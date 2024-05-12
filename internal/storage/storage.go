package storage

import (
	"fmt"
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
	if _, ok := storage.CounterItems[key]; ok {
		storage.CounterItems[key] = storage.CounterItems[key] + value
	} else {
		storage.CounterItems[key] = value
	}
	return true
}

func (storage *MemStorage) GetGaugeItem(key string) (model.Gauge, bool) {
	if _, ok := storage.GaugeItems[key]; ok {
		return storage.GaugeItems[key], true
	}
	return 0, false
}

func (storage *MemStorage) GetCounterItem(key string) (model.Counter, bool) {
	if _, ok := storage.CounterItems[key]; ok {
		return storage.CounterItems[key], true
	}
	return 0, false
}

func (storage *MemStorage) GetAllGaugeItems() (map[string]string, bool) {
	res := map[string]string{}
	for key, value := range storage.GaugeItems {
		res[key] = fmt.Sprintf("%.2f", value)
	}
	return res, true
}

func (storage *MemStorage) GetAllCounterItems() (map[string]string, bool) {
	res := map[string]string{}
	for key, value := range storage.CounterItems {
		res[key] = fmt.Sprintf("%v", value)
	}
	return res, true
}
