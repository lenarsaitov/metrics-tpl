package routers

import (
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lenarsaitov/metrics-tpl/internal/server/config"
	"github.com/lenarsaitov/metrics-tpl/internal/server/controllers"
	"github.com/lenarsaitov/metrics-tpl/internal/server/repository/inmemory"
	"github.com/lenarsaitov/metrics-tpl/internal/server/repository/postgres"
	"github.com/lenarsaitov/metrics-tpl/internal/server/routers/middlewares"
	"github.com/lenarsaitov/metrics-tpl/internal/server/services"
	"github.com/lenarsaitov/metrics-tpl/internal/server/worker"
	"net/http"
)

func GetRouters(ctx context.Context, cfg *config.Config) (*echo.Echo, error) {
	e := echo.New()

	var storage services.Storage
	var err error
	if cfg.DatabaseDSN != "" {
		storage, err = postgres.NewPollStorage(ctx, cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}
	} else {
		storage = inmemory.NewPollStorage()
	}

	e.Use(
		middlewares.ApplyRequestInform,
		middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}),
	)

	useMetrics := services.NewMetricsService(storage)
	worker.New(useMetrics, cfg.StoreInterval, cfg.FileStoragePath).Run(ctx, cfg.Restore)
	serverController := controllers.New(cfg.DatabaseDSN, useMetrics)

	e.Add(http.MethodGet, "/ping", serverController.PingDB)

	e.Add(http.MethodPost, "/value/", serverController.GetMetric)
	e.Add(http.MethodPost, "/update/", serverController.Update)

	e.Add(http.MethodPost, "/updates/", serverController.Updates)

	e.Add(http.MethodGet, "/value/:metricType/:metricName", serverController.GetMetricPath)
	e.Add(http.MethodPost, "/update/:metricType/:metricName/:metricValue", serverController.UpdatePath)

	e.Add(http.MethodGet, "/", serverController.GetAllMetrics)

	return e, nil
}
