package implementations

import (
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"math/rand"
	"runtime"
)

type MetricPollModel struct {
	pollCount int64
}

var _ models.MetricPoll = &MetricPollModel{}

func NewMetricPollModel() *MetricPollModel {
	return &MetricPollModel{}
}

func (m *MetricPollModel) GetAgentMetrics() models.Metrics {
	var memStats = &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	metrics := matchMemStatsAgentMetrics(memStats)
	metrics = append(
		metrics,
		models.Metric{MetricName: PollCountMetric, MetricType: models.CounterMetricType, MetricValue: float64(m.pollCount)},
		models.Metric{MetricName: RandomValueMetric, MetricType: models.GaugeMetricType, MetricValue: rand.Float64()},
	)

	return metrics
}

func matchMemStatsAgentMetrics(m *runtime.MemStats) models.Metrics {
	return models.Metrics{
		models.Metric{MetricName: AllocMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.Alloc)},
		models.Metric{MetricName: BuckHashSysMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.BuckHashSys)},
		models.Metric{MetricName: FreesMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.Frees)},
		models.Metric{MetricName: GCCPUFractionMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.GCCPUFraction)},
		models.Metric{MetricName: GCSysMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.GCSys)},
		models.Metric{MetricName: HeapAllocMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.HeapAlloc)},
		models.Metric{MetricName: HeapIdleMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.HeapIdle)},
		models.Metric{MetricName: HeapInuseMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.HeapInuse)},
		models.Metric{MetricName: HeapObjectsMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.HeapObjects)},
		models.Metric{MetricName: HeapReleasedMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.HeapReleased)},
		models.Metric{MetricName: HeapSysMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.HeapSys)},
		models.Metric{MetricName: LastGCMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.LastGC)},
		models.Metric{MetricName: LookupsMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.Lookups)},
		models.Metric{MetricName: MCacheInuseMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.MCacheInuse)},
		models.Metric{MetricName: MCacheSysMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.MCacheSys)},
		models.Metric{MetricName: MSpanInuseMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.MSpanInuse)},
		models.Metric{MetricName: MSpanSysMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.MSpanSys)},
		models.Metric{MetricName: MallocsMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.Mallocs)},
		models.Metric{MetricName: NextGCMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.NextGC)},
		models.Metric{MetricName: NumForcedGCMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.NumForcedGC)},
		models.Metric{MetricName: NumGCMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.NumGC)},
		models.Metric{MetricName: OtherSysMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.OtherSys)},
		models.Metric{MetricName: PauseTotalNsMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.PauseTotalNs)},
		models.Metric{MetricName: StackInuseMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.StackInuse)},
		models.Metric{MetricName: StackSysMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.StackSys)},
		models.Metric{MetricName: SysMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.Sys)},
		models.Metric{MetricName: TotalAllocMetric, MetricType: models.GaugeMetricType, MetricValue: float64(m.TotalAlloc)},
	}
}
