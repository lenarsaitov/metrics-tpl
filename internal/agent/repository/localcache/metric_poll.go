package localcache

import (
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"math/rand"
	"runtime"
)

type MetricPollStorage struct {
	pollCount int64
}

func NewMetricPollStorage() *MetricPollStorage {
	return &MetricPollStorage{}
}

func (m *MetricPollStorage) GetAgentMetrics() models.Metrics {
	var memStats = &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	metrics := matchMemStatsAgentMetrics(memStats)
	metrics.CounterMetrics = append(metrics.CounterMetrics, models.CounterMetric{Name: models.PollCountMetric, Value: m.pollCount})
	metrics.GaugeMetrics = append(metrics.GaugeMetrics, models.GaugeMetric{Name: models.RandomValueMetric, Value: rand.Float64()})

	return *metrics
}

func matchMemStatsAgentMetrics(m *runtime.MemStats) *models.Metrics {
	gaugeMetrics := []models.GaugeMetric{
		{Name: models.AllocMetric, Value: float64(m.Alloc)},
		{Name: models.BuckHashSysMetric, Value: float64(m.BuckHashSys)},
		{Name: models.FreesMetric, Value: float64(m.Frees)},
		{Name: models.GCCPUFractionMetric, Value: float64(m.GCCPUFraction)},
		{Name: models.GCSysMetric, Value: float64(m.GCSys)},
		{Name: models.HeapAllocMetric, Value: float64(m.HeapAlloc)},
		{Name: models.HeapIdleMetric, Value: float64(m.HeapIdle)},
		{Name: models.HeapInuseMetric, Value: float64(m.HeapInuse)},
		{Name: models.HeapObjectsMetric, Value: float64(m.HeapObjects)},
		{Name: models.HeapReleasedMetric, Value: float64(m.HeapReleased)},
		{Name: models.HeapSysMetric, Value: float64(m.HeapSys)},
		{Name: models.LastGCMetric, Value: float64(m.LastGC)},
		{Name: models.LookupsMetric, Value: float64(m.Lookups)},
		{Name: models.MCacheInuseMetric, Value: float64(m.MCacheInuse)},
		{Name: models.MCacheSysMetric, Value: float64(m.MCacheSys)},
		{Name: models.MSpanInuseMetric, Value: float64(m.MSpanInuse)},
		{Name: models.MSpanSysMetric, Value: float64(m.MSpanSys)},
		{Name: models.MallocsMetric, Value: float64(m.Mallocs)},
		{Name: models.NextGCMetric, Value: float64(m.NextGC)},
		{Name: models.NumForcedGCMetric, Value: float64(m.NumForcedGC)},
		{Name: models.NumGCMetric, Value: float64(m.NumGC)},
		{Name: models.OtherSysMetric, Value: float64(m.OtherSys)},
		{Name: models.PauseTotalNsMetric, Value: float64(m.PauseTotalNs)},
		{Name: models.StackInuseMetric, Value: float64(m.StackInuse)},
		{Name: models.StackSysMetric, Value: float64(m.StackSys)},
		{Name: models.SysMetric, Value: float64(m.Sys)},
		{Name: models.TotalAllocMetric, Value: float64(m.TotalAlloc)},
	}

	return &models.Metrics{GaugeMetrics: gaugeMetrics}
}
