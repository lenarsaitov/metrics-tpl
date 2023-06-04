package inmemory

import (
	"context"
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

func (m *PollStorage) GetAll(ctx context.Context) (models.Metrics, error) {
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

	return metrics, nil
}

func (m *PollStorage) GetGaugeMetric(ctx context.Context, name string) (*float64, error) {
	if value, ok := m.gaugeMetrics[name]; ok {
		return &value, nil
	}

	return nil, nil
}

func (m *PollStorage) GetCounterMetric(ctx context.Context, name string) (*int64, error) {
	if value, ok := m.counterMetrics[name]; ok {
		return &value, nil
	}

	return nil, nil
}

func (m *PollStorage) ReplaceGauge(ctx context.Context, name string, value float64) error {
	m.gaugeMetrics[name] = value

	return nil
}

func (m *PollStorage) AddCounter(ctx context.Context, name string, value int64) (int64, error) {
	m.counterMetrics[name] = m.counterMetrics[name] + value

	return m.counterMetrics[name], nil
}
