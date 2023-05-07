package memstorage

const (
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

type Service interface {
	ReplaceGauge(name string, value float64)
	AddCounter(name string, value int64)
}
