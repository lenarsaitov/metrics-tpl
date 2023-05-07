package models

const (
	CounterMetricType = "counter"
	GaugeMetricType   = "gauge"
)

type Metric struct {
	MetricType  string
	MetricName  string
	MetricValue float64
}

type Metrics []Metric

type MetricPoll interface {
	GetAgentMetrics() Metrics
}
