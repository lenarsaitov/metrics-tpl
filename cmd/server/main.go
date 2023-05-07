package main

import (
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/controllers"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models/implementations"
	"github.com/lenarsaitov/metrics-tpl/internal/server/usecase"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("start metrics collection and alerting service web server..")

	parseConfiguration()

	e := echo.New()
	useMetrics := usecase.NewMetricsUseCase(implementations.NewMemStorageModel())
	serverController := controllers.New(useMetrics)

	e.Add(http.MethodGet, "/", serverController.GetMetrics)
	e.Add(http.MethodGet, "/value/:metricType/:metricName", serverController.GetMetric)
	e.Add(http.MethodPost, "/update/:metricType/:metricName/:metricValue", serverController.Update)

	log.Info().Msg("running server on: " + flagAddrRun)
	if err := http.ListenAndServe(flagAddrRun, e); err != nil {
		log.Fatal().Err(err).Msg("failed to run web server")
	}
}
