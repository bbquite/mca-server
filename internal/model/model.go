package model

type Gauge float64
type Counter int64

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsPack struct {
	Metrics []Metric `json:"metrics"`
}

// CREATE TYPE metric_type AS ENUM (
//     'GAUGE',
//     'COUNTER'
// );
// create table metrics (
// 	id serial PRIMARY KEY,
// 	metric_type metric_type not null,
// 	metric_name varchar(55) UNIQUE not null,
//     delta integer,
//     value double precision
// );
