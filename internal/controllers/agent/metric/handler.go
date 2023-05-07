package metric

import (
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/metriclisten"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/metricsender"
	"github.com/rs/zerolog"
	"sync"
	"time"
)

type Handler struct {
	metricListenService metriclisten.Service
	metricSenderService metricsender.Service
	listenedMetrics     metriclisten.Metrics
	pollCount           int64
	pollInterval        int
	reportInterval      int
}

func NewHandler(
	metricListenService metriclisten.Service,
	metricSenderService metricsender.Service,
	pollInterval int,
	reportInterval int,
) *Handler {
	return &Handler{
		metricListenService: metricListenService,
		metricSenderService: metricSenderService,
		listenedMetrics:     make(metriclisten.Metrics, 0),
		pollInterval:        pollInterval,
		reportInterval:      reportInterval,
	}
}

func (h *Handler) Handle(log *zerolog.Logger) {
	var mu = &sync.Mutex{}

	log.Info().Msg("start get metrics...")
	go h.getMetrics(mu)

	time.Sleep(time.Second * time.Duration(h.pollInterval))

	log.Info().Msg("start send metrics...")
	for {
		for _, metric := range h.listenedMetrics {
			switch metric.MetricType {
			case metriclisten.CounterMetricType:
				err := h.metricSenderService.SendAddCounter(metric.MetricName, int64(metric.MetricValue))
				if err != nil {
					log.Error().Err(err).Msg("failed to send counter metric")

					return
				}
			case metriclisten.GaugeMetricType:
				err := h.metricSenderService.SendReplaceGauge(metric.MetricName, metric.MetricValue)
				if err != nil {
					log.Error().Err(err).Msg("failed to send gauge metric")

					return
				}
			default:
				log.Warn().Str("metric_type", metric.MetricType).Msg("invalid metric type")
			}
		}

		log.Info().Interface("metrics", h.listenedMetrics).Msg("sent metrics")
		time.Sleep(time.Second * time.Duration(h.reportInterval))
	}
}

func (h *Handler) getMetrics(mu *sync.Mutex) {
	for {
		metrics := h.metricListenService.GetMetrics()

		mu.Lock()
		h.listenedMetrics = metrics
		mu.Unlock()

		time.Sleep(time.Second * time.Duration(h.pollInterval))
	}
}
