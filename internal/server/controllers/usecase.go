package controllers

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/rs/zerolog"
)

type (
	MetricsServerUseCase interface {
		GetAllMetrics() models.Metrics
		GetMetric(metricType, metricName string) *float64
		UpdateGaugeMetric(log *zerolog.Logger, metricName string, metricValue string) error
		UpdateCounterMetric(log *zerolog.Logger, metricName string, metricValue string) error
	}
)
