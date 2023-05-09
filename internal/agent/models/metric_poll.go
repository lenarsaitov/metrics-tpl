package models

const (
	CounterMetricType = "counter"
	GaugeMetricType   = "gauge"
)

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

type Metrics struct {
	GaugeMetrics   []GaugeMetric
	CounterMetrics []CounterMetric
}

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
