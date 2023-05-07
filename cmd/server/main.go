package main

import (
	"github.com/labstack/echo"
	"net/http"

	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server"
	"github.com/lenarsaitov/metrics-tpl/internal/models/implementations"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("start metrics collection and alerting service web server..")

	e := echo.New()
	serverController := server.NewController(implementations.NewMemStorageModel())

	e.Add(http.MethodGet, "/", serverController.GetMetrics)
	e.Add(http.MethodGet, "/value/:metricType/:metricName", serverController.GetMetric)
	e.Add(http.MethodPost, "/update/:metricType/:metricName/:metricValue", serverController.Update)

	if err := http.ListenAndServe("localhost:8080", e); err != nil {
		log.Fatal().Err(err).Msg("failed to run web server")
	}
}
