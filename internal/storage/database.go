package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/bbquite/mca-server/internal/model"
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
			if pgErr.Code == "42P01" { // relation does not exist
				_, err = storage.DB.ExecContext(storage.ctx, sqlCreateString)
				if err != nil {
					log.Print(err)
					return err
				}
			}
		}
		log.Print(err)
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

func (storage *DBStorage) AddMetricItem(mType string, key string, value any) bool {
	// storage.mx.Lock()
	// defer storage.mx.Unlock()

	// sqlString := `
	// 	EXISTS (SELECT id FROM metrics WHERE metric_name = $2)
	// 	BEGIN
	// 		UPDATE metrics
	// 		SET value = $3
	// 		WHERE metric_name = $2
	// 	END
	// 	ELSE BEGIN
	// 		INSERT INTO metrics (metric_type, metric_name, value)
	// 		VALUES ($1, $2, $3)
	// 	END
	// `

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
		log.Print(err)
		return false
	}

	return true

	// sqlString := `
	// 	INSERT INTO metrics (metric_type, metric_name, value)
	// 	SELECT $1, $2, $3
	// 	WHERE
	// 	NOT EXISTS (
	// 		SELECT id FROM metrics WHERE metric_name = $2
	// 	)
	// `

	// var metricID uint8

	// sqlStringSelect := `
	// 	SELECT id
	// 	FROM metrics
	// 	WHERE metric_name = $1
	// 	LIMIT 1
	// `

	// sqlStringInsert := `
	// 	INSERT INTO metrics (metric_type, metric_name, value)
	// 	VALUES ($1, $2, $3)
	// `

	// sqlStringUpdate := `
	// 	UPDATE metrics
	// 	SET value = $2
	// 	WHERE metric_name = $1
	// `

	// if mType == "COUNTER" {
	// 	sqlStringUpdate = `
	// 		UPDATE metrics
	// 		SET delta = delta + $2
	// 		WHERE metric_name = $1
	// 	`
	// 	sqlStringInsert = `
	// 		INSERT INTO metrics (metric_type, metric_name, delta)
	// 		VALUES ($1, $2, $3)
	// 	`
	// }

	// row := storage.DB.QueryRowContext(storage.ctx, sqlStringSelect, key)
	// err := row.Scan(&metricID)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		_, err = storage.DB.ExecContext(storage.ctx, sqlStringInsert, mType, key, value)
	// 		if err != nil {
	// 			log.Print(err)
	// 			return false
	// 		}
	// 		return true
	// 	}
	// 	log.Print(err)
	// 	return false
	// }

	// _, err = storage.DB.ExecContext(storage.ctx, sqlStringUpdate, key, value)
	// if err != nil {
	// 	log.Print(err)
	// 	return false
	// }

	// return true
}

func (storage *DBStorage) AddGaugeItem(key string, value model.Gauge) bool {
	return storage.AddMetricItem("GAUGE", key, value)
}

func (storage *DBStorage) AddCounterItem(key string, value model.Counter) bool {
	return storage.AddMetricItem("COUNTER", key, value)
}

func (storage *DBStorage) ResetCounterItem(key string) bool {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	var metricID uint8

	sqlStringSelect := `
		SELECT id 
		FROM metrics 
		WHERE metric_name = $1 
		LIMIT 1
	`

	sqlStringUpdate := `
		UPDATE metrics
		SET delta = delta + $2
		WHERE metric_name = $1
	`

	row := storage.DB.QueryRowContext(storage.ctx, sqlStringSelect, key)
	err := row.Scan(&metricID)
	if err != nil {
		log.Print(err)
		return false
	}

	_, err = storage.DB.ExecContext(storage.ctx, sqlStringUpdate, key, 0)
	if err != nil {
		log.Print(err)
		return false
	}

	return false
}

func (storage *DBStorage) GetGaugeItem(key string) (model.Gauge, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	var metric model.Gauge = 0

	sqlStringSelect := `
		SELECT value 
		FROM metrics 
		WHERE metric_name = $1 
		LIMIT 1
	`

	row := storage.DB.QueryRowContext(storage.ctx, sqlStringSelect, key)
	err := row.Scan(&metric)
	if err != nil {
		log.Print(err)
		return metric, false
	}
	return metric, true
}

func (storage *DBStorage) GetCounterItem(key string) (model.Counter, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()

	var metric model.Counter = 0

	sqlStringSelect := `
		SELECT delta 
		FROM metrics 
		WHERE metric_name = $1 
		LIMIT 1
	`

	row := storage.DB.QueryRowContext(storage.ctx, sqlStringSelect, key)
	err := row.Scan(&metric)
	if err != nil {
		log.Print(err)
		return metric, false
	}
	return metric, false
}

func (storage *DBStorage) GetGaugeItems() (map[string]model.Gauge, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()
	// result := storage.GaugeItems
	result := make(map[string]model.Gauge)

	sqlStringSelect := `
		SELECT metric_name, value 
		FROM metrics 
		WHERE metric_type = 'GAUGE'
	`
	rows, err := storage.DB.QueryContext(storage.ctx, sqlStringSelect)
	if err != nil {
		log.Print(err)
		return nil, false
	}
	for rows.Next() {
		var metricName string
		var metricValue model.Gauge

		err := rows.Scan(&metricName, &metricValue)
		if err != nil {
			log.Print(err)
			return nil, false
		}

		result[metricName] = metricValue
	}

	return result, true
}

func (storage *DBStorage) GetCounterItems() (map[string]model.Counter, bool) {
	// storage.mx.RLock()
	// defer storage.mx.RUnlock()
	// result := storage.CounterItems
	result := make(map[string]model.Counter)

	sqlStringSelect := `
		SELECT metric_name, delta 
		FROM metrics 
		WHERE metric_type = 'COUNTER'
	`
	rows, err := storage.DB.QueryContext(storage.ctx, sqlStringSelect)
	if err != nil {
		log.Print(err)
		return nil, false
	}
	for rows.Next() {
		var metricName string
		var metricValue model.Counter

		err := rows.Scan(&metricName, &metricValue)
		if err != nil {
			log.Print(err)
			return nil, false
		}

		result[metricName] = metricValue
	}

	return result, true
}

func (storage *DBStorage) AddMetricsPack(metrics *model.MetricsPack) bool {

	tx, err := storage.DB.Begin()
	if err != nil {
		log.Print(err)
		return false
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
			log.Print(err)
			return false
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return false
	}

	return true

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
