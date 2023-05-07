package implementations

import (
	"github.com/lenarsaitov/metrics-tpl/internal/models/services"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/memstorage"
)

type MemStorageModel struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string]int64
}

var _ memstorage.Service = &MemStorageModel{}

func NewMemStorageModel() *MemStorageModel {
	return &MemStorageModel{
		gaugeMetrics:   map[string]float64{},
		counterMetrics: map[string]int64{},
	}
}

func (m *MemStorageModel) GetAllMetrics() memstorage.ServerMetrics {
	metrics := memstorage.ServerMetrics{}

	for name, value := range m.gaugeMetrics {
		metrics = append(metrics, memstorage.ServerMetric{
			MetricType:  services.GaugeMetricType,
			MetricName:  name,
			MetricValue: value,
		})
	}

	for name, value := range m.counterMetrics {
		metrics = append(metrics, memstorage.ServerMetric{
			MetricType:  services.CounterMetricType,
			MetricName:  name,
			MetricValue: float64(value),
		})
	}

	return metrics
}

func (m *MemStorageModel) GetGaugeMetric(name string) *float64 {
	if value, ok := m.gaugeMetrics[name]; ok {
		return &value
	}

	return nil
}

func (m *MemStorageModel) GetCounterMetric(name string) *int64 {
	if value, ok := m.counterMetrics[name]; ok {
		return &value
	}

	return nil
}

func (m *MemStorageModel) ReplaceGauge(name string, value float64) {
	m.gaugeMetrics[name] = value
}

func (m *MemStorageModel) AddCounter(name string, value int64) {
	m.counterMetrics[name] = m.counterMetrics[name] + value
}
