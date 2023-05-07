package controllers

import (
	"github.com/google/uuid"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/usecase"
	logger "github.com/rs/zerolog/log"
)

type Controller struct {
	metricsUseCase usecase.MetricsAgent
}

func New(metricsUseCase usecase.MetricsAgent) *Controller {
	return &Controller{
		metricsUseCase: metricsUseCase,
	}
}

func (c *Controller) ListenAndSend() {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	c.metricsUseCase.PollAndReport(&log)
}
