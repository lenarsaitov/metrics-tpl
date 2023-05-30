package routers

import (
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/controllers"
	"github.com/lenarsaitov/metrics-tpl/internal/server/middlewares"
	"github.com/lenarsaitov/metrics-tpl/internal/server/repository"
	"github.com/lenarsaitov/metrics-tpl/internal/server/services"
	"net/http"
)

func GetRouters() *echo.Echo {
	e := echo.New()
	useMetrics := services.NewMetricsService(repository.NewPollStorage())
	serverController := controllers.New(useMetrics)

	//e.Use(middlewares.ApplyRequestInform, middlewares.ApplyGZIP)
	e.Use(middlewares.ApplyRequestInform)

	e.Add(http.MethodPost, "/value/", serverController.GetMetric)
	e.Add(http.MethodPost, "/update/", serverController.Update)

	e.Add(http.MethodGet, "/value/:metricType/:metricName", serverController.GetMetricPath)
	e.Add(http.MethodPost, "/update/:metricType/:metricName/:metricValue", serverController.UpdatePath)

	e.Add(http.MethodGet, "/", serverController.GetAllMetrics)

	return e
}
