package main

import (
	"context"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/config"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/repository"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/services"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os/signal"
	"syscall"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("start metrics agent for collecting runtime metrics and then report them to the server via HTTP protocol..")

	cfg, err := config.GetConfiguration()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get configuration, flag or environment")
	}

	log.Info().
		Str("server_address_remote", cfg.RemoteAddr).
		Int("poll_interval", cfg.PollInterval).
		Int("report_interval", cfg.ReportInterval).
		Msg("agent settings")

	ctx, cnl := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cnl()

	service := services.NewMetricsService(
		repository.NewPollStorage(),
		cfg.RemoteAddr,
		cfg.PollInterval,
		cfg.ReportInterval,
	)

	agentWorker := worker.New(service)
	agentWorker.Run(ctx)
}
