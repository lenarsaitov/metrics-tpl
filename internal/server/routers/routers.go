package routers

import (
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/controllers"
	"github.com/lenarsaitov/metrics-tpl/internal/server/repository"
	"github.com/lenarsaitov/metrics-tpl/internal/server/services"
	"net/http"
)

func GetRouters() *echo.Echo {
	e := echo.New()
	useMetrics := services.NewMetricsService(repository.NewPollStorage())
	serverController := controllers.New(useMetrics)

	e.Add(http.MethodGet, "/", serverController.GetAllMetrics)
	e.Add(http.MethodGet, "/value/:metricType/:metricName", serverController.GetMetric)
	e.Add(http.MethodPost, "/update/:metricType/:metricName/:metricValue", serverController.Update)

	return e
}
