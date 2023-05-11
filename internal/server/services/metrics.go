package services

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/rs/zerolog"
	"strconv"
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

func (h *MetricsService) UpdateGaugeMetric(log *zerolog.Logger, metricName string, metricValue string) error {
	gaugeValue, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return err
	}

	h.storage.ReplaceGauge(metricName, gaugeValue)

	log.Info().
		Str("metric_name", metricName).
		Str("metric_value", metricValue).
		Msg("gauge was replaced successfully")

	return nil
}

func (h *MetricsService) UpdateCounterMetric(log *zerolog.Logger, metricName string, metricValue string) error {
	countValue, err := strconv.Atoi(metricValue)
	if err != nil {
		return err
	}

	h.storage.AddCounter(metricName, int64(countValue))

	log.Info().
		Str("metric_name", metricName).
		Str("metric_value", metricValue).
		Msg("counter was added successfully")

	return nil
}
