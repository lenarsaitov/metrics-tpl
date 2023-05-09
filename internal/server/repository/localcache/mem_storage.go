package localcache

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
)

type MemStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMetrics:   map[string]float64{},
		counterMetrics: map[string]int64{},
	}
}

func (m *MemStorage) GetAllMetrics() models.Metrics {
	metrics := models.Metrics{}

	for name, value := range m.gaugeMetrics {
		metrics = append(metrics, models.Metric{
			MetricType:  models.GaugeMetricType,
			MetricName:  name,
			MetricValue: value,
		})
	}

	for name, value := range m.counterMetrics {
		metrics = append(metrics, models.Metric{
			MetricType:  models.CounterMetricType,
			MetricName:  name,
			MetricValue: float64(value),
		})
	}

	return metrics
}

func (m *MemStorage) GetGaugeMetric(name string) *float64 {
	if value, ok := m.gaugeMetrics[name]; ok {
		return &value
	}

	return nil
}

func (m *MemStorage) GetCounterMetric(name string) *int64 {
	if value, ok := m.counterMetrics[name]; ok {
		return &value
	}

	return nil
}

func (m *MemStorage) ReplaceGauge(name string, value float64) {
	m.gaugeMetrics[name] = value
}

func (m *MemStorage) AddCounter(name string, value int64) {
	m.counterMetrics[name] = m.counterMetrics[name] + value
}
