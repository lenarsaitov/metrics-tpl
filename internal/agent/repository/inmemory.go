package repository

import (
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"math/rand"
	"runtime"
)

const (
	AllocMetric         = "Alloc"
	BuckHashSysMetric   = "BuckHashSys"
	FreesMetric         = "Frees"
	GCCPUFractionMetric = "GCCPUFraction"
	GCSysMetric         = "GCSys"
	HeapAllocMetric     = "HeapAlloc"
	HeapIdleMetric      = "HeapIdle"
	HeapInuseMetric     = "HeapInuse"
	HeapObjectsMetric   = "HeapObjects"
	HeapReleasedMetric  = "HeapReleased"
	HeapSysMetric       = "HeapSys"
	LastGCMetric        = "LastGC"
	LookupsMetric       = "Lookups"
	MCacheInuseMetric   = "MCacheInuse"
	MCacheSysMetric     = "MCacheSys"
	MSpanInuseMetric    = "MSpanInuse"
	MSpanSysMetric      = "MSpanSys"
	MallocsMetric       = "Mallocs"
	NextGCMetric        = "NextGC"
	NumForcedGCMetric   = "NumForcedGC"
	NumGCMetric         = "NumGC"
	OtherSysMetric      = "OtherSys"
	PauseTotalNsMetric  = "PauseTotalNs"
	StackInuseMetric    = "StackInuse"
	StackSysMetric      = "StackSys"
	SysMetric           = "Sys"
	TotalAllocMetric    = "TotalAlloc"

	// PollCountMetric и RandomValueMetric метрики не из runtime
	PollCountMetric   = "PollCount"
	RandomValueMetric = "RandomValue"
)

type PollStorage struct {
	pollCount int64
}

func NewPollStorage() *PollStorage {
	return &PollStorage{}
}

func (m *PollStorage) GetPoll() models.Metrics {
	metrics := getRuntimePoll()
	m.pollCount++

	metrics.CounterMetrics = append(metrics.CounterMetrics, models.CounterMetric{Name: PollCountMetric, Value: m.pollCount})
	metrics.GaugeMetrics = append(metrics.GaugeMetrics, models.GaugeMetric{Name: RandomValueMetric, Value: rand.Float64()})

	return *metrics
}

func getRuntimePoll() *models.Metrics {
	var m = &runtime.MemStats{}
	runtime.ReadMemStats(m)

	gaugeMetrics := []models.GaugeMetric{
		{Name: AllocMetric, Value: float64(m.Alloc)},
		{Name: BuckHashSysMetric, Value: float64(m.BuckHashSys)},
		{Name: FreesMetric, Value: float64(m.Frees)},
		{Name: GCCPUFractionMetric, Value: float64(m.GCCPUFraction)},
		{Name: GCSysMetric, Value: float64(m.GCSys)},
		{Name: HeapAllocMetric, Value: float64(m.HeapAlloc)},
		{Name: HeapIdleMetric, Value: float64(m.HeapIdle)},
		{Name: HeapInuseMetric, Value: float64(m.HeapInuse)},
		{Name: HeapObjectsMetric, Value: float64(m.HeapObjects)},
		{Name: HeapReleasedMetric, Value: float64(m.HeapReleased)},
		{Name: HeapSysMetric, Value: float64(m.HeapSys)},
		{Name: LastGCMetric, Value: float64(m.LastGC)},
		{Name: LookupsMetric, Value: float64(m.Lookups)},
		{Name: MCacheInuseMetric, Value: float64(m.MCacheInuse)},
		{Name: MCacheSysMetric, Value: float64(m.MCacheSys)},
		{Name: MSpanInuseMetric, Value: float64(m.MSpanInuse)},
		{Name: MSpanSysMetric, Value: float64(m.MSpanSys)},
		{Name: MallocsMetric, Value: float64(m.Mallocs)},
		{Name: NextGCMetric, Value: float64(m.NextGC)},
		{Name: NumForcedGCMetric, Value: float64(m.NumForcedGC)},
		{Name: NumGCMetric, Value: float64(m.NumGC)},
		{Name: OtherSysMetric, Value: float64(m.OtherSys)},
		{Name: PauseTotalNsMetric, Value: float64(m.PauseTotalNs)},
		{Name: StackInuseMetric, Value: float64(m.StackInuse)},
		{Name: StackSysMetric, Value: float64(m.StackSys)},
		{Name: SysMetric, Value: float64(m.Sys)},
		{Name: TotalAllocMetric, Value: float64(m.TotalAlloc)},
	}

	return &models.Metrics{GaugeMetrics: gaugeMetrics}
}
