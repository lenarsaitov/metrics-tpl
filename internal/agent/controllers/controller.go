package controllers

import (
	"github.com/google/uuid"
	logger "github.com/rs/zerolog/log"
)

type Controller struct {
	metricsUseCase MetricsAgent
}

func New(metricsUseCase MetricsAgent) *Controller {
	return &Controller{
		metricsUseCase: metricsUseCase,
	}
}

func (c *Controller) ListenAndSend() {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	c.metricsUseCase.PollAndReport(&log)
}
