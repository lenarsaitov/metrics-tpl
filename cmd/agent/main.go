package main

import (
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/agent"
	"github.com/lenarsaitov/metrics-tpl/internal/models/implementations"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("start metrics agent for collecting runtime metrics and then sending them to the server via HTTP protocol..")

	agentController := agent.NewController(
		implementations.NewMetricListenModel(),
		implementations.NewMetricSenderModel(),
		2,
		10,
	)

	agentController.ListenAndSend()
}
