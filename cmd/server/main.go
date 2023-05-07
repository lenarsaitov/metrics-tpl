package main

import (
	"net/http"

	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server"
	"github.com/lenarsaitov/metrics-tpl/internal/models/implementations"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("start metrics collection and alerting service web server..")

	serverController := server.NewController(implementations.NewMemStorageModel())

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", serverController.Update)

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		log.Fatal().Err(err).Msg("failed to run web server")
	}
}
