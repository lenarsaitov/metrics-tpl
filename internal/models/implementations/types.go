package implementations

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

const (
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

type ServerResponse struct {
	Response struct {
		Text string `json:"text,omitempty"`
	} `json:"response,omitempty"`
}
