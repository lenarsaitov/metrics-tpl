package worker

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

type (
	MetricsService interface {
		PollWithTicker(ctx context.Context, log *zerolog.Logger)
		ReportWithTicker(ctx context.Context, log *zerolog.Logger)
	}
)

type Worker struct {
	metricsService MetricsService
}

func New(metricsService MetricsService) *Worker {
	return &Worker{
		metricsService: metricsService,
	}
}

func (r *Worker) Run(ctx context.Context) {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	log.Info().Msg("start poll metrics...")
	go r.metricsService.PollWithTicker(ctx, &log)

	log.Info().Msg("start report metrics...")
	r.metricsService.ReportWithTicker(ctx, &log)
}
