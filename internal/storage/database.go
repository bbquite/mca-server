package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/bbquite/mca-server/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	DB  *sql.DB
	ctx context.Context
}

func NewDBStorage(ctx context.Context, databaseDSN string) (*DBStorage, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return &DBStorage{}, err
	}

	return &DBStorage{
		DB:  db,
		ctx: ctx,
	}, nil
}

func (storage *DBStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := storage.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (storage *DBStorage) AddGaugeItem(key string, value model.Gauge) bool {
	// storage.mx.Lock()
	// defer storage.mx.Unlock()
	// storage.GaugeItems[key] = value
	return true
}

func (storage *DBStorage) AddCounterItem(key string, value model.Counter) bool {
	// storage.mx.Lock()
	// defer storage.mx.Unlock()

	// if _, ok := storage.CounterItems[key]; ok {
	// 	storage.CounterItems[key] = storage.CounterItems[key] + value
	// } else {
	// 	storage.CounterItems[key] = value
	// }
	return true
}

func (storage *DBStorage) GetGaugeItem(key string) (model.Gauge, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	// if _, ok := storage.GaugeItems[key]; ok {
	// 	return storage.GaugeItems[key], true
	// }
	return 0, false
}

func (storage *DBStorage) GetCounterItem(key string) (model.Counter, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	// if _, ok := storage.CounterItems[key]; ok {
	// 	return storage.CounterItems[key], true
	// }
	return 0, false
}

func (storage *DBStorage) ResetCounterItem(key string) bool {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	// if _, ok := storage.CounterItems[key]; ok {
	// 	storage.CounterItems[key] = model.Counter(0)
	// 	return true
	// }
	return false
}

func (storage *DBStorage) GetGaugeItems() (map[string]model.Gauge, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()
	// result := storage.GaugeItems
	result := make(map[string]model.Gauge)
	return result, true
}

func (storage *DBStorage) GetCounterItems() (map[string]model.Counter, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()
	// result := storage.CounterItems
	result := make(map[string]model.Counter)
	return result, true
}

func (storage *DBStorage) GetStringGaugeItems() (map[string]string, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	res := map[string]string{}
	// for key, value := range storage.GaugeItems {
	// 	res[key] = fmt.Sprintf("%.2f", value)
	// }
	return res, true
}

func (storage *DBStorage) GetStringCounterItems() (map[string]string, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	res := map[string]string{}
	// for key, value := range storage.CounterItems {
	// 	res[key] = fmt.Sprintf("%v", value)
	// }
	return res, true
}
