package controllers

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

type (
	MetricsService interface {
		PollAndReport(log *zerolog.Logger)
	}
)

type Controller struct {
	metricsService MetricsService
}

func New(metricsService MetricsService) *Controller {
	return &Controller{
		metricsService: metricsService,
	}
}

func (c *Controller) PollAndReport() {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	c.metricsService.PollAndReport(&log)
}
