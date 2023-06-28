package worker

import (
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

type (
	MetricsService interface {
		PutCommonPollWorker(ctx context.Context, log *zerolog.Logger)
		PutPsutilPollWorker(ctx context.Context, log *zerolog.Logger)
		WriteToChanWorker(ctx context.Context, log *zerolog.Logger)
		SendWorker(ctx context.Context, log *zerolog.Logger)
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

func (r *Worker) Run(ctx context.Context, rateLimit int) {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	log.Info().Msg("start put poll common metrics...")
	go r.metricsService.PutCommonPollWorker(ctx, &log)

	log.Info().Msg("start put poll psutil metrics...")
	go r.metricsService.PutPsutilPollWorker(ctx, &log)

	log.Info().Msg("start report to in jobs metrics...")
	go r.metricsService.WriteToChanWorker(ctx, &log)

	log.Info().Int("count", rateLimit).Msg("start jobs with reporting metrics...")
	for i := 0; i < rateLimit; i++ {
		go r.metricsService.SendWorker(ctx, &log)
	}
}
