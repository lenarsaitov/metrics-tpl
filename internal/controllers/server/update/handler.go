package update

import (
	"github.com/lenarsaitov/metrics-tpl/internal/models/services"
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

func (h *Handler) Handle(log *zerolog.Logger, rsp *responder.Responder, input *Input) error {
	log.Info().Msg("handle update request")

	switch input.MetricType {
	case services.GaugeMetricType:
		gaugeValue, err := strconv.ParseFloat(input.MetricValue, 64)
		if err != nil {
			return rsp.BadRequest("invalid value of gauge metrics, need float64")
		}

		h.memStorageService.ReplaceGauge(input.MetricName, gaugeValue)

		log.Info().
			Str("metric_name", input.MetricName).
			Str("metric_value", input.MetricValue).
			Msg("gauge was replaced successfully")

		return rsp.OK("gauge was replaced successfully")
	case services.CounterMetricType:
		countValue, err := strconv.Atoi(input.MetricValue)
		if err != nil {
			return rsp.BadRequest("invalid value of counter metrics, need int64")
		}

		h.memStorageService.AddCounter(input.MetricName, int64(countValue))

		log.Info().
			Str("metric_name", input.MetricName).
			Str("metric_value", input.MetricValue).
			Msg("counter was added successfully")

		return rsp.OK("counter was added successfully")
	default:
		return rsp.BadRequest("unavailable metric type, use counter or gauge")
	}
}
