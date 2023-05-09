package main

import (
	"github.com/lenarsaitov/metrics-tpl/internal/server/config"
	"github.com/lenarsaitov/metrics-tpl/internal/server/routers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("start metrics collection and alerting service web server..")

	cfg, err := config.GetConfiguration()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get configuration, flag or environment")
	}

	log.Info().Msg("running server on: " + cfg.AddrRun)
	if err := http.ListenAndServe(cfg.AddrRun, routers.GetRouters()); err != nil {
		log.Fatal().Err(err).Msg("failed to run web server")
	}
}
