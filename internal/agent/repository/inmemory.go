package repository

import (
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
	"runtime"
	"sync"
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

	// PollCountMetric и RandomValueMetric метрики типа counter, не из runtime
	PollCountMetric   = "PollCount"
	RandomValueMetric = "RandomValue"

	// Метрики типа gauge из psutil
	TotalMemory     = "TotalMemory"
	FreeMemory      = "FreeMemory"
	CPUutilization1 = "CPUutilization1"
)

type PollStorage struct {
	mx        *sync.RWMutex
	metrics   models.Metrics
	pollCount int64
}

func NewPollStorage() *PollStorage {
	return &PollStorage{
		mx: new(sync.RWMutex),
	}
}

func (m *PollStorage) PutCommonPoll() {
	metrics := getRuntimePoll()

	m.mx.Lock()
	defer m.mx.Unlock()

	m.pollCount++

	metrics.CounterMetrics = append(metrics.CounterMetrics, models.CounterMetric{Name: PollCountMetric, Value: m.pollCount})
	metrics.GaugeMetrics = append(metrics.GaugeMetrics, models.GaugeMetric{Name: RandomValueMetric, Value: rand.Float64()})

	m.metrics.CounterMetrics = append(m.metrics.CounterMetrics, metrics.CounterMetrics...)
	m.metrics.GaugeMetrics = append(m.metrics.GaugeMetrics, metrics.GaugeMetrics...)
}

func (m *PollStorage) PutPsutilPoll() {
	v, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	m.mx.Lock()
	defer m.mx.Unlock()

	m.metrics.GaugeMetrics = append(m.metrics.GaugeMetrics, models.GaugeMetric{Name: TotalMemory, Value: float64(v.Total)})
	m.metrics.GaugeMetrics = append(m.metrics.GaugeMetrics, models.GaugeMetric{Name: FreeMemory, Value: float64(v.Free)})
	//m.metrics.GaugeMetrics = append(m.metrics.GaugeMetrics, models.GaugeMetric{Name: CPUutilization1, Value: float64(v.)})
}

func (m *PollStorage) GetPoll() models.Metrics {
	m.mx.Lock()
	defer m.mx.Unlock()

	return m.metrics
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
