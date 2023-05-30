package services

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
)

type Storage interface {
	GetAllMetrics() models.Metrics

	GetGaugeMetric(name string) *float64
	GetCounterMetric(name string) *int64

	ReplaceGauge(name string, value float64)
	AddCounter(name string, value int64)
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

func (h *MetricsService) GetMetric(metricType, metricName string) *float64 {
	switch metricType {
	case models.GaugeMetricType:
		value := h.storage.GetGaugeMetric(metricName)
		if value == nil {
			return nil
		}

		return value
	case models.CounterMetricType:
		value := h.storage.GetCounterMetric(metricName)
		if value == nil {
			return nil
		}

		valueFloat64 := float64(*value)

		return &valueFloat64
	}

	return nil
}

func (h *MetricsService) UpdateGaugeMetric(metricName string, gaugeValue float64) error {
	h.storage.ReplaceGauge(metricName, gaugeValue)

	return nil
}

func (h *MetricsService) UpdateCounterMetric(metricName string, counterValue int64) error {
	h.storage.AddCounter(metricName, counterValue)

	return nil
}
