package main

import (
	"github.com/lenarsaitov/metrics-tpl/internal/agent/controllers"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models/implementations"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/usecase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("start metrics agent for collecting runtime metrics and then report them to the server via HTTP protocol..")

	if err := parseConfiguration(); err != nil {
		log.Fatal().Err(err).Msg("failed to parse configuration, flag or environment")
	}

	log.Info().
		Str("server_address_remote", flagRemoteAddr).
		Int("poll_interval", flagPollInterval).
		Int("report_interval", flagReportInterval).
		Msg("agent settings")

	useMetrics := usecase.NewMetricsUseCase(
		implementations.NewMetricPollModel(),
		implementations.NewMetricReportModel(flagRemoteAddr),
		flagPollInterval,
		flagReportInterval,
	)

	agentController := controllers.New(useMetrics)

	agentController.ListenAndSend()
}
