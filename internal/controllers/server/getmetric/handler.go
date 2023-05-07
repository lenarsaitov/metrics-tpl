package getmetric

import (
	"github.com/lenarsaitov/metrics-tpl/internal/models/services"
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
	log.Info().Msg("handle get request")

	switch input.MetricType {
	case services.GaugeMetricType:
		value := h.memStorageService.GetGaugeMetric(input.MetricName)
		if value == nil {
			return rsp.NotFound("not found value of this gauge metric")
		}

		return rsp.OKWithBody(*value)
	case services.CounterMetricType:
		values := h.memStorageService.GetCounterMetric(input.MetricName)
		if values == nil {
			return rsp.NotFound("not found values of this counter metric")
		}

		return rsp.OKWithBody((*values)[len(*values)-1])
	default:
		return rsp.BadRequest("unavailable metric type, use counter or gauge")
	}
}
