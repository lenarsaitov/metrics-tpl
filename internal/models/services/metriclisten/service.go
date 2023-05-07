package metriclisten

type Metric struct {
	MetricType  string
	MetricName  string
	MetricValue float64
}

type Metrics []Metric

const (
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

type Service interface {
	GetMetrics() Metrics
}
