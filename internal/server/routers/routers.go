package routers

import (
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lenarsaitov/metrics-tpl/internal/server/config"
	"github.com/lenarsaitov/metrics-tpl/internal/server/controllers"
	"github.com/lenarsaitov/metrics-tpl/internal/server/repository/inmemory"
	"github.com/lenarsaitov/metrics-tpl/internal/server/routers/middlewares"
	"github.com/lenarsaitov/metrics-tpl/internal/server/runner"
	"github.com/lenarsaitov/metrics-tpl/internal/server/services"
	"net/http"
)

func GetRouters(cfg *config.Config) *echo.Echo {
	e := echo.New()

	memoryStorage := inmemory.NewPollStorage()
	useMetrics := services.NewMetricsService(memoryStorage)

	runner.New(useMetrics, cfg.StoreInterval, cfg.FileStoragePath).Run(context.Background(), cfg.Restore)
	serverController := controllers.New(cfg.DatabaseDSN, useMetrics)

	e.Use(middlewares.ApplyRequestInform, middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Add(http.MethodGet, "/ping", serverController.PingDB)

	e.Add(http.MethodPost, "/value/", serverController.GetMetric)
	e.Add(http.MethodPost, "/update/", serverController.Update)

	e.Add(http.MethodGet, "/value/:metricType/:metricName", serverController.GetMetricPath)
	e.Add(http.MethodPost, "/update/:metricType/:metricName/:metricValue", serverController.UpdatePath)

	e.Add(http.MethodGet, "/", serverController.GetAllMetrics)

	return e
}
