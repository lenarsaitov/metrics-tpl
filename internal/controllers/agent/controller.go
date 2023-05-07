package agent

import (
	"github.com/google/uuid"
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/agent/metric"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/metriclisten"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/metricsender"
	logger "github.com/rs/zerolog/log"
)

type Controller struct {
	metricHandler *metric.Handler
}

func NewController(
	metricListenService metriclisten.Service,
	metricSenderService metricsender.Service,
	pollInterval int,
	reportInterval int,
) *Controller {
	return &Controller{
		metricHandler: metric.NewHandler(metricListenService, metricSenderService, pollInterval, reportInterval),
	}
}

func (c *Controller) ListenAndSend() {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	c.metricHandler.Handle(
		&log,
	)
}
