package implementations

import "github.com/lenarsaitov/metrics-tpl/internal/models/services/memstorage"

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

func (m *MemStorageModel) ReplaceGauge(name string, value float64) {
	m.gaugeMetrics[name] = value
}

func (m *MemStorageModel) AddCounter(name string, value int64) {
	if _, ok := m.counterMetrics[name]; !ok {
		m.counterMetrics[name] = make([]int64, 0)
	}

	m.counterMetrics[name] = append(m.counterMetrics[name], value)
}
