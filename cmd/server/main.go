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

	cfg := config.GetConfiguration()

	log.Info().Msg("running server on: " + cfg.AddrRun)
	if err := http.ListenAndServe(cfg.AddrRun, routers.GetRouters(cfg)); err != nil {
		log.Fatal().Err(err).Msg("failed to run web server")
	}
}
