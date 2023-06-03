package worker

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

type (
	MetricsService interface {
		Poll(ctx context.Context, log *zerolog.Logger)
		Report(ctx context.Context, log *zerolog.Logger)
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
	go r.metricsService.Poll(ctx, &log)

	log.Info().Msg("start report metrics...")
	r.metricsService.Report(ctx, &log)
}
