package models

const (
	CounterMetricType = "counter"
	GaugeMetricType   = "gauge"
)

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

type Metrics struct {
	GaugeMetrics   []GaugeMetric
	CounterMetrics []CounterMetric
}
