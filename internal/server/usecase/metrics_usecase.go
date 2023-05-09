package usecase

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/rs/zerolog"
	"strconv"
)

type MetricsUseCase struct {
	memStorage MemStorage
}

func NewMetricsUseCase(memStorageService MemStorage) *MetricsUseCase {
	return &MetricsUseCase{
		memStorage: memStorageService,
	}
}

func (h *MetricsUseCase) GetAllMetrics() models.Metrics {
	return h.memStorage.GetAllMetrics()
}

func (h *MetricsUseCase) GetMetric(metricType, metricName string) *float64 {
	switch metricType {
	case models.GaugeMetricType:
		value := h.memStorage.GetGaugeMetric(metricName)
		if value == nil {
			return nil
		}

		return value
	case models.CounterMetricType:
		value := h.memStorage.GetCounterMetric(metricName)
		if value == nil {
			return nil
		}

		valueFloat64 := float64(*value)
		return &valueFloat64
	}

	return nil
}

func (h *MetricsUseCase) UpdateGaugeMetric(log *zerolog.Logger, metricName string, metricValue string) error {
	gaugeValue, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return err
	}

	h.memStorage.ReplaceGauge(metricName, gaugeValue)

	log.Info().
		Str("metric_name", metricName).
		Str("metric_value", metricValue).
		Msg("gauge was replaced successfully")

	return nil
}

func (h *MetricsUseCase) UpdateCounterMetric(log *zerolog.Logger, metricName string, metricValue string) error {
	countValue, err := strconv.Atoi(metricValue)
	if err != nil {
		return err
	}

	h.memStorage.AddCounter(metricName, int64(countValue))

	log.Info().
		Str("metric_name", metricName).
		Str("metric_value", metricValue).
		Msg("counter was added successfully")

	return nil
}
