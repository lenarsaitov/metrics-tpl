package usecase

import "github.com/lenarsaitov/metrics-tpl/internal/agent/models"

type MetricPoll interface {
	GetAgentMetrics() models.Metrics
}

type MetricReport interface {
	ReportGaugeMetric(name string, value float64) error
	ReportCounterMetric(name string, value int64) error
}
