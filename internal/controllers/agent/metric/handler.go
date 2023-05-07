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
	listenedMetrics     []*metriclisten.Metrics
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
		pollInterval:        pollInterval,
		reportInterval:      reportInterval,
	}
}

func (h *Handler) Handle(log *zerolog.Logger) {
	var mu = &sync.Mutex{}

	log.Info().Msg("start send metrics...")
	go func(log *zerolog.Logger, mu *sync.Mutex) {
		for {
			if len(h.listenedMetrics) == 0 {
				continue
			}

			mu.Lock()
			lastMetrics := *h.listenedMetrics[len(h.listenedMetrics)-1]
			mu.Unlock()
			log.Info().Interface("last_metrics", lastMetrics).Msg("send last metrics to server")

			for _, metric := range lastMetrics {
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
			mu.Lock()
			h.listenedMetrics = h.listenedMetrics[:len(h.listenedMetrics)-1]
			mu.Unlock()
			time.Sleep(time.Second * time.Duration(h.reportInterval))
		}
	}(log, mu)

	log.Info().Msg("start get metrics...")
	for {
		metrics := h.metricListenService.GetMetrics()
		log.Info().Interface("metrics", metrics).Msg("got metrics")
		h.listenedMetrics = append(h.listenedMetrics, metrics)
		time.Sleep(time.Second * time.Duration(h.pollInterval))
	}
}
