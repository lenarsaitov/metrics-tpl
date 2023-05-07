package implementations

import (
	"github.com/lenarsaitov/metrics-tpl/internal/models/services"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/memstorage"
)

type MemStorageModel struct {
	gaugeMetrics   map[string]float64
	counterMetrics map[string][]int64
}

var _ memstorage.Service = &MemStorageModel{}

func NewMemStorageModel() *MemStorageModel {
	return &MemStorageModel{
		gaugeMetrics:   map[string]float64{},
		counterMetrics: map[string][]int64{},
	}
}

func (m *MemStorageModel) GetAllMetrics() memstorage.ServerMetrics {
	metrics := memstorage.ServerMetrics{}

	for name, value := range m.gaugeMetrics {
		metrics = append(metrics, memstorage.ServerMetric{
			MetricType:  services.GaugeMetricType,
			MetricName:  name,
			MetricValue: []float64{value},
		})
	}

	for name, value := range m.counterMetrics {
		metrics = append(metrics, memstorage.ServerMetric{
			MetricType:  services.CounterMetricType,
			MetricName:  name,
			MetricValue: transformSliceIntToFloat(value),
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

func (m *MemStorageModel) GetCounterMetric(name string) *[]int64 {
	if values, ok := m.counterMetrics[name]; ok && len(values) != 0 {
		return &values
	}

	return nil
}

func (m *MemStorageModel) ReplaceGauge(name string, value float64) {
	m.gaugeMetrics[name] = value
}

func (m *MemStorageModel) AddCounter(name string, value int64) {
	if _, ok := m.counterMetrics[name]; !ok {
		m.counterMetrics[name] = make([]int64, 0)
	}

	m.counterMetrics[name] = append(m.counterMetrics[name], value)
}

func transformSliceIntToFloat(numbersInt []int64) []float64 {
	numbersFloat := make([]float64, 0, len(numbersInt))

	for _, number := range numbersInt {
		numbersFloat = append(numbersFloat, float64(number))
	}

	return numbersFloat
}
