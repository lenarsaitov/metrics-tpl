package runner

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

type (
	MetricsService interface {
		PollAndReport(ctx context.Context, log *zerolog.Logger)
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

func (r *Runner) Run(ctx context.Context) {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	r.metricsService.PollAndReport(ctx, &log)
}
