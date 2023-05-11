package runner

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

type Runner struct {
	metricsService MetricsService
}

func New(metricsService MetricsService) *Runner {
	return &Runner{
		metricsService: metricsService,
	}
}

func (c *Runner) PollAndReport() {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	c.metricsService.PollAndReport(&log)
}
