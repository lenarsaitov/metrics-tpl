package getmetrics

import (
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

func (h *Handler) Handle(log *zerolog.Logger, rsp *responder.Responder) error {
	log.Info().Msg("handle get request")

	return rsp.OKWithBody(h.memStorageService.GetAllMetrics())
}
