package usecase

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/rs/zerolog"
)

type (
	MetricsAgent interface {
		PollAndReport(log *zerolog.Logger)
	}

	MetricsServer interface {
		GetAllMetrics() models.Metrics
		GetMetric(metricType, metricName string) *float64
		UpdateGaugeMetric(log *zerolog.Logger, metricName string, metricValue string) error
		UpdateCounterMetric(log *zerolog.Logger, metricName string, metricValue string) error
	}
)
