package update

import (
	"fmt"
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

	fmt.Println(input.MetricType)

	switch input.MetricType {
	case memstorage.GaugeMetricType:
		gaugeValue, err := strconv.ParseFloat(input.MetricValue, 64)
		if err != nil {
			log.Error().Err(err).Msg("invalid value of gauge metrics, need float64")
			rsp.InternalError()

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
			log.Error().Err(err).Msg("invalid value of counter metrics, need int64")
			rsp.InternalError()

			return
		}

		h.memStorageService.AddCounter(input.MetricName, int64(countValue))

		log.Info().
			Str("metric_name", input.MetricName).
			Str("metric_value", input.MetricValue).
			Msg("counter was added successfully")

		rsp.OK("counter was added successfully")
	default:
		log.Info().Str("metric_type", input.MetricName).Msg("unavailable metric type")

		rsp.InternalError()
	}

	log.Info().Msg("handled")
}
