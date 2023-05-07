package metricsender

type Service interface {
	SendReplaceGauge(name string, value float64) error
	SendAddCounter(name string, value int64) error
}
