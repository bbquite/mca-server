package storage

import "errors"

var (
	ErrorAddingGauge   = errors.New("no gauge value added")
	ErrorAddingCounter = errors.New("no counter value added")

	ErrorGaugeNotFound   = errors.New("gauge not found")
	ErrorCounterNotFound = errors.New("counters not found")
	ErrorGettingMetrics  = errors.New("error getting metrics")

	ErrorResetCounter = errors.New("error reset counter")
)
