package implementations

import (
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

func (m *MetricListenModel) GetMetrics() metriclisten.Metrics {
	var memStats = &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	metrics := matchMemStatsMetrics(memStats)
	metrics = append(
		metrics,
		metriclisten.Metric{MetricName: PollCountMetric, MetricType: metriclisten.CounterMetricType, MetricValue: float64(m.pollCount)},
		metriclisten.Metric{MetricName: RandomValueMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: rand.Float64()},
	)

	return metrics
}

func matchMemStatsMetrics(m *runtime.MemStats) metriclisten.Metrics {
	return metriclisten.Metrics{
		metriclisten.Metric{MetricName: AllocMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.Alloc)},
		metriclisten.Metric{MetricName: BuckHashSysMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.BuckHashSys)},
		metriclisten.Metric{MetricName: FreesMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.Frees)},
		metriclisten.Metric{MetricName: GCCPUFractionMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.GCCPUFraction)},
		metriclisten.Metric{MetricName: GCSysMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.GCSys)},
		metriclisten.Metric{MetricName: HeapAllocMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.HeapAlloc)},
		metriclisten.Metric{MetricName: HeapIdleMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.HeapIdle)},
		metriclisten.Metric{MetricName: HeapInuseMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.HeapInuse)},
		metriclisten.Metric{MetricName: HeapObjectsMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.HeapObjects)},
		metriclisten.Metric{MetricName: HeapReleasedMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.HeapReleased)},
		metriclisten.Metric{MetricName: HeapSysMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.HeapSys)},
		metriclisten.Metric{MetricName: LastGCMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.LastGC)},
		metriclisten.Metric{MetricName: LookupsMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.Lookups)},
		metriclisten.Metric{MetricName: MCacheInuseMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.MCacheInuse)},
		metriclisten.Metric{MetricName: MCacheSysMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.MCacheSys)},
		metriclisten.Metric{MetricName: MSpanInuseMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.MSpanInuse)},
		metriclisten.Metric{MetricName: MSpanSysMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.MSpanSys)},
		metriclisten.Metric{MetricName: MallocsMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.Mallocs)},
		metriclisten.Metric{MetricName: NextGCMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.NextGC)},
		metriclisten.Metric{MetricName: NumForcedGCMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.NumForcedGC)},
		metriclisten.Metric{MetricName: NumGCMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.NumGC)},
		metriclisten.Metric{MetricName: OtherSysMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.OtherSys)},
		metriclisten.Metric{MetricName: PauseTotalNsMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.PauseTotalNs)},
		metriclisten.Metric{MetricName: StackInuseMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.StackInuse)},
		metriclisten.Metric{MetricName: StackSysMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.StackSys)},
		metriclisten.Metric{MetricName: SysMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.Sys)},
		metriclisten.Metric{MetricName: TotalAllocMetric, MetricType: metriclisten.GaugeMetricType, MetricValue: float64(m.TotalAlloc)},
	}
}
