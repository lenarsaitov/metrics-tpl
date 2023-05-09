package usecase

import "github.com/lenarsaitov/metrics-tpl/internal/server/models"

type MemStorage interface {
	GetAllMetrics() models.Metrics

	GetGaugeMetric(name string) *float64
	GetCounterMetric(name string) *int64

	ReplaceGauge(name string, value float64)
	AddCounter(name string, value int64)
}
