package routers

import (
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/controllers"
	"github.com/lenarsaitov/metrics-tpl/internal/server/repository/localcache"
	"github.com/lenarsaitov/metrics-tpl/internal/server/usecase"
	"net/http"
)

func GetRouters() *echo.Echo {
	e := echo.New()
	useMetrics := usecase.NewMetricsUseCase(localcache.NewMemStorage())
	serverController := controllers.New(useMetrics)

	e.Add(http.MethodGet, "/", serverController.GetMetrics)
	e.Add(http.MethodGet, "/value/:metricType/:metricName", serverController.GetMetric)
	e.Add(http.MethodPost, "/update/:metricType/:metricName/:metricValue", serverController.Update)

	return e
}
