package main

import (
	"github.com/lenarsaitov/metrics-tpl/internal/agent/config"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/controllers"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/repository/localcache"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/usecase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	useMetrics := usecase.NewMetricsUseCase(
		localcache.NewMetricPollStorage(),
		localcache.NewMetricReportStorage(cfg.RemoteAddr),
		cfg.PollInterval,
		cfg.ReportInterval,
	)

	agentController := controllers.New(useMetrics)

	agentController.ListenAndSend()
}
