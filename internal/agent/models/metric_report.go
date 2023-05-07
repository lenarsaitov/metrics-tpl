package models

type MetricReport interface {
	ReportGaugeMetric(name string, value float64) error
	ReportCounterMetric(name string, value int64) error
}
