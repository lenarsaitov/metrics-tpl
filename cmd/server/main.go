package main

import (
	"context"
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
		log.Fatal().Err(err).Msg("failed to get configuration")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e, err := routers.GetRouters(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get routers")
	}

	log.Info().Msg("running server on: " + cfg.AddrRun)
	if err := http.ListenAndServe(cfg.AddrRun, e); err != nil {
		log.Fatal().Err(err).Msg("failed to run web server")
	}
}
