package usecase

import (
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"github.com/rs/zerolog"
	"sync"
	"time"
)

type MetricsUseCase struct {
	metricPollService   MetricPoll
	metricReportService MetricReport
	polledMetrics       models.Metrics
	pollCount           int64
	pollInterval        int
	reportInterval      int
}

func NewMetricsUseCase(
	metricPollService MetricPoll,
	metricReportService MetricReport,
	pollInterval int,
	reportInterval int,
) *MetricsUseCase {
	return &MetricsUseCase{
		metricPollService:   metricPollService,
		metricReportService: metricReportService,
		polledMetrics:       models.Metrics{GaugeMetrics: make([]models.GaugeMetric, 0), CounterMetrics: make([]models.CounterMetric, 0)},
		pollInterval:        pollInterval,
		reportInterval:      reportInterval,
	}
}

func (h *MetricsUseCase) PollAndReport(log *zerolog.Logger) {
	var mu = &sync.Mutex{}

	log.Info().Msg("start poll metrics...")
	go h.getMetrics(mu)

	time.Sleep(time.Second * time.Duration(h.pollInterval))

	log.Info().Msg("start report metrics...")
	for {
		for _, metric := range h.polledMetrics.GaugeMetrics {
			err := h.metricReportService.ReportGaugeMetric(metric.Name, metric.Value)
			if err != nil {
				log.Error().Err(err).Msg("failed to report gauge metric")

				return
			}
		}
		for _, metric := range h.polledMetrics.CounterMetrics {
			err := h.metricReportService.ReportCounterMetric(metric.Name, metric.Value)
			if err != nil {
				log.Error().Err(err).Msg("failed to report counter metric")

				return
			}
		}

		log.Info().Interface("metrics", h.polledMetrics).Msg("reported metrics")
		time.Sleep(time.Second * time.Duration(h.reportInterval))
	}
}

func (h *MetricsUseCase) getMetrics(mu *sync.Mutex) {
	for {
		metrics := h.metricPollService.GetAgentMetrics()

		mu.Lock()
		h.polledMetrics = metrics
		mu.Unlock()

		time.Sleep(time.Second * time.Duration(h.pollInterval))
	}
}
