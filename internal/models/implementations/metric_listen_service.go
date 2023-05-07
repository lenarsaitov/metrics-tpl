package implementations

import (
	"github.com/lenarsaitov/metrics-tpl/internal/models/services"
	"math/rand"
	"runtime"

	"github.com/lenarsaitov/metrics-tpl/internal/models/services/metriclisten"
)

type MetricListenModel struct {
	pollCount int64
}

var _ metriclisten.Service = &MetricListenModel{}

func NewMetricListenModel() *MetricListenModel {
	return &MetricListenModel{}
}

func (m *MetricListenModel) GetAgentMetrics() metriclisten.AgentMetrics {
	var memStats = &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	metrics := matchMemStatsAgentMetrics(memStats)
	metrics = append(
		metrics,
		metriclisten.AgentMetric{MetricName: PollCountMetric, MetricType: services.CounterMetricType, MetricValue: float64(m.pollCount)},
		metriclisten.AgentMetric{MetricName: RandomValueMetric, MetricType: services.GaugeMetricType, MetricValue: rand.Float64()},
	)

	return metrics
}

func matchMemStatsAgentMetrics(m *runtime.MemStats) metriclisten.AgentMetrics {
	return metriclisten.AgentMetrics{
		metriclisten.AgentMetric{MetricName: AllocMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.Alloc)},
		metriclisten.AgentMetric{MetricName: BuckHashSysMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.BuckHashSys)},
		metriclisten.AgentMetric{MetricName: FreesMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.Frees)},
		metriclisten.AgentMetric{MetricName: GCCPUFractionMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.GCCPUFraction)},
		metriclisten.AgentMetric{MetricName: GCSysMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.GCSys)},
		metriclisten.AgentMetric{MetricName: HeapAllocMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.HeapAlloc)},
		metriclisten.AgentMetric{MetricName: HeapIdleMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.HeapIdle)},
		metriclisten.AgentMetric{MetricName: HeapInuseMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.HeapInuse)},
		metriclisten.AgentMetric{MetricName: HeapObjectsMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.HeapObjects)},
		metriclisten.AgentMetric{MetricName: HeapReleasedMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.HeapReleased)},
		metriclisten.AgentMetric{MetricName: HeapSysMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.HeapSys)},
		metriclisten.AgentMetric{MetricName: LastGCMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.LastGC)},
		metriclisten.AgentMetric{MetricName: LookupsMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.Lookups)},
		metriclisten.AgentMetric{MetricName: MCacheInuseMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.MCacheInuse)},
		metriclisten.AgentMetric{MetricName: MCacheSysMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.MCacheSys)},
		metriclisten.AgentMetric{MetricName: MSpanInuseMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.MSpanInuse)},
		metriclisten.AgentMetric{MetricName: MSpanSysMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.MSpanSys)},
		metriclisten.AgentMetric{MetricName: MallocsMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.Mallocs)},
		metriclisten.AgentMetric{MetricName: NextGCMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.NextGC)},
		metriclisten.AgentMetric{MetricName: NumForcedGCMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.NumForcedGC)},
		metriclisten.AgentMetric{MetricName: NumGCMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.NumGC)},
		metriclisten.AgentMetric{MetricName: OtherSysMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.OtherSys)},
		metriclisten.AgentMetric{MetricName: PauseTotalNsMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.PauseTotalNs)},
		metriclisten.AgentMetric{MetricName: StackInuseMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.StackInuse)},
		metriclisten.AgentMetric{MetricName: StackSysMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.StackSys)},
		metriclisten.AgentMetric{MetricName: SysMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.Sys)},
		metriclisten.AgentMetric{MetricName: TotalAllocMetric, MetricType: services.GaugeMetricType, MetricValue: float64(m.TotalAlloc)},
	}
}
