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

type MemStorage interface {
	GetAllMetrics() Metrics

	GetGaugeMetric(name string) *float64
	GetCounterMetric(name string) *int64

	ReplaceGauge(name string, value float64)
	AddCounter(name string, value int64)
}
