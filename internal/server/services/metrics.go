package services

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
)

type Storage interface {
	GetAllMetrics() models.Metrics

	GetGaugeMetric(name string) *float64
	GetCounterMetric(name string) *int64

	ReplaceGauge(name string, value float64)
	AddCounter(name string, value int64) int64
}

type MetricsService struct {
	storage Storage
}

func NewMetricsService(storageService Storage) *MetricsService {
	return &MetricsService{
		storage: storageService,
	}
}

func (h *MetricsService) GetAllMetrics() models.Metrics {
	return h.storage.GetAllMetrics()
}

func (h *MetricsService) GetGaugeMetric(metricName string) *float64 {
	return h.storage.GetGaugeMetric(metricName)
}

func (h *MetricsService) GetCounterMetric(metricName string) *int64 {
	return h.storage.GetCounterMetric(metricName)
}

func (h *MetricsService) UpdateGaugeMetric(metricName string, gaugeValue float64) {
	h.storage.ReplaceGauge(metricName, gaugeValue)
}

func (h *MetricsService) UpdateCounterMetric(metricName string, counterValue int64) int64 {
	return h.storage.AddCounter(metricName, counterValue)
}
