package controllers

import "github.com/lenarsaitov/metrics-tpl/internal/server/models"

type (
	MetricsService interface {
		GetAllMetrics() models.Metrics
		GetGaugeMetric(metricName string) *float64
		GetCounterMetric(metricName string) *int64
		UpdateGaugeMetric(metricName string, gaugeValue float64)
		UpdateCounterMetric(metricName string, counterValue int64) int64
	}
)
