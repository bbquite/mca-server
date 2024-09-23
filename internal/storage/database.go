package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/bbquite/mca-server/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (storage *DBStorage) CheckDatabaseValid() error {
	err := storage.Ping()
	if err != nil {
		return err
	}

	sqlCheckString := `SELECT id FROM metrics LIMIT 1`
	sqlCreateString := `
		DROP TYPE IF EXISTS metric_type;
		CREATE TYPE metric_type AS ENUM('GAUGE','COUNTER');

		create table metrics (
			id serial PRIMARY KEY,
			metric_type metric_type not null,
			metric_name varchar(55) UNIQUE not null,
			delta integer,
			value double precision
		);
	`

	_, err = storage.DB.ExecContext(storage.ctx, sqlCheckString)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UndefinedTable {
				_, err = storage.DB.ExecContext(storage.ctx, sqlCreateString)
				if err != nil {
					return err
				}
			}
		}
		return err
	}

	return nil
}

func (storage *DBStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := storage.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (storage *DBStorage) AddMetricItem(mType string, key string, value any) error {
	// storage.mx.Lock()
	// defer storage.mx.Unlock()

	sqlString := `
		INSERT INTO metrics (metric_type, metric_name, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (metric_name) DO UPDATE SET value = $3
	`

	if mType == "COUNTER" {
		sqlString = `
			INSERT INTO metrics (metric_type, metric_name, delta)
			VALUES ($1, $2, $3)
			ON CONFLICT (metric_name) DO UPDATE SET delta = metrics.delta + $3
		`
	}

	_, err := storage.DB.ExecContext(storage.ctx, sqlString, mType, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) AddGaugeItem(key string, value model.Gauge) error {
	return storage.AddMetricItem("GAUGE", key, value)
}

func (storage *DBStorage) AddCounterItem(key string, value model.Counter) error {
	return storage.AddMetricItem("COUNTER", key, value)
}

func (storage *DBStorage) GetGaugeItem(key string) (model.Gauge, error) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	var metric model.Gauge

	sqlStringSelect := `
		SELECT value 
		FROM metrics 
		WHERE metric_name = $1 
		LIMIT 1
	`

	row := storage.DB.QueryRowContext(storage.ctx, sqlStringSelect, key)
	err := row.Scan(&metric)
	if err != nil {
		return 0, err
	}
	return metric, nil
}

func (storage *DBStorage) GetCounterItem(key string) (model.Counter, error) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	var metric model.Counter

	sqlStringSelect := `
		SELECT delta 
		FROM metrics 
		WHERE metric_name = $1 
		LIMIT 1
	`

	row := storage.DB.QueryRowContext(storage.ctx, sqlStringSelect, key)
	err := row.Scan(&metric)
	if err != nil {
		return 0, err
	}
	return metric, nil
}

func (storage *DBStorage) GetGaugeItems() (map[string]model.Gauge, error) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	result := make(map[string]model.Gauge)

	sqlStringSelect := `
		SELECT metric_name, value 
		FROM metrics 
		WHERE metric_type = 'GAUGE'
	`

	rows, err := storage.DB.QueryContext(storage.ctx, sqlStringSelect)
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, err
	}

	for rows.Next() {
		var metricName string
		var metricValue model.Gauge

		err := rows.Scan(&metricName, &metricValue)
		if err != nil {
			return nil, err
		}

		result[metricName] = metricValue
	}

	return result, nil
}

func (storage *DBStorage) GetCounterItems() (map[string]model.Counter, error) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	result := make(map[string]model.Counter)

	sqlStringSelect := `
		SELECT metric_name, delta 
		FROM metrics 
		WHERE metric_type = 'COUNTER'
	`

	rows, err := storage.DB.QueryContext(storage.ctx, sqlStringSelect)
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, err
	}

	for rows.Next() {
		var metricName string
		var metricValue model.Counter

		err := rows.Scan(&metricName, &metricValue)
		if err != nil {
			return nil, err
		}

		result[metricName] = metricValue
	}

	return result, nil
}

func (storage *DBStorage) AddMetricsPack(metrics *model.MetricsPack) error {

	tx, err := storage.DB.Begin()
	if err != nil {
		return err
	}

	var sqlString string = ""
	var value any

	for _, el := range metrics.Metrics {
		mType := strings.ToUpper(el.MType)
		switch mType {
		case "GAUGE":
			sqlString = `
				INSERT INTO metrics (metric_type, metric_name, value)
				VALUES ($1, $2, $3)
				ON CONFLICT (metric_name) DO UPDATE SET value = $3
			`
			value = el.Value
		case "COUNTER":
			sqlString = `
				INSERT INTO metrics (metric_type, metric_name, delta)
				VALUES ($1, $2, $3)
				ON CONFLICT (metric_name) DO UPDATE SET delta = metrics.delta + $3
			`
			value = el.Delta
		}

		_, err := tx.ExecContext(storage.ctx, sqlString, mType, el.ID, value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (storage *DBStorage) ResetCounterItem(key string) error {
	return errors.New("UNUSED")
}
