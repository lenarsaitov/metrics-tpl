package memstorage

type ServerMetric struct {
	MetricType  string
	MetricName  string
	MetricValue []float64
}

type ServerMetrics []ServerMetric

type Service interface {
	GetAllMetrics() ServerMetrics

	GetGaugeMetric(name string) *float64
	GetCounterMetric(name string) *[]int64

	ReplaceGauge(name string, value float64)
	AddCounter(name string, value int64)
}
