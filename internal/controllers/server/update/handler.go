package update

import (
	"strconv"

	"github.com/lenarsaitov/metrics-tpl/internal/models/services/memstorage"
	"github.com/lenarsaitov/metrics-tpl/internal/responder"
	"github.com/rs/zerolog"
)

type Handler struct {
	memStorageService memstorage.Service
}

func NewHandler(memStorageService memstorage.Service) *Handler {
	return &Handler{
		memStorageService: memStorageService,
	}
}

func (h *Handler) Handle(log *zerolog.Logger, rsp *responder.Responder, input *Input) {
	log.Info().Msg("handle")

	switch input.MetricType {
	case memstorage.GaugeMetricType:
		gaugeValue, err := strconv.ParseFloat(input.MetricValue, 64)
		if err != nil {
			rsp.BadRequest("invalid value of gauge metrics, need float64")

			return
		}

		h.memStorageService.ReplaceGauge(input.MetricName, gaugeValue)

		log.Info().
			Str("metric_name", input.MetricName).
			Str("metric_value", input.MetricValue).
			Msg("gauge was replaced successfully")

		rsp.OK("gauge was replaced successfully")
	case memstorage.CounterMetricType:
		countValue, err := strconv.Atoi(input.MetricValue)
		if err != nil {
			rsp.BadRequest("invalid value of counter metrics, need int64")

			return
		}

		h.memStorageService.AddCounter(input.MetricName, int64(countValue))

		log.Info().
			Str("metric_name", input.MetricName).
			Str("metric_value", input.MetricValue).
			Msg("counter was added successfully")

		rsp.OK("counter was added successfully")
	default:
		rsp.BadRequest("unavailable metric type, use counter or gauge")
	}

	log.Info().Msg("handled")
}
