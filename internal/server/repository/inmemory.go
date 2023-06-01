package repository

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
)

type PollStorage struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

func NewPollStorage() *PollStorage {
	return &PollStorage{
		gaugeMetrics:   map[string]float64{},
		counterMetrics: map[string]int64{},
	}
}

func (m *PollStorage) GetAll() models.Metrics {
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

func (m *PollStorage) GetGaugeMetric(name string) *float64 {
	if value, ok := m.gaugeMetrics[name]; ok {
		return &value
	}

	return nil
}

func (m *PollStorage) GetCounterMetric(name string) *int64 {
	if value, ok := m.counterMetrics[name]; ok {
		return &value
	}

	return nil
}

func (m *PollStorage) ReplaceGauge(name string, value float64) {
	m.gaugeMetrics[name] = value
}

func (m *PollStorage) AddCounter(name string, value int64) int64 {
	m.counterMetrics[name] = m.counterMetrics[name] + value

	return m.counterMetrics[name]
}
