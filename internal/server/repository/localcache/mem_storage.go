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
	metrics := models.Metrics{
		GaugeMetrics:   make([]models.GaugeMetric, 0, len(m.gaugeMetrics)),
		CounterMetrics: make([]models.CounterMetric, 0, len(m.counterMetrics)),
	}

	for name, value := range m.gaugeMetrics {
		metrics.GaugeMetrics = append(metrics.GaugeMetrics, models.GaugeMetric{
			Name:  name,
			Value: value,
		})
	}

	for name, value := range m.counterMetrics {
		metrics.CounterMetrics = append(metrics.CounterMetrics, models.CounterMetric{
			Name:  name,
			Value: value,
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
