package controllers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"net/http"
)

type (
	MetricsService interface {
		GetAllMetrics() models.Metrics
		GetMetric(metricType, metricName string) *float64
		UpdateGaugeMetric(log *zerolog.Logger, metricName string, metricValue string) error
		UpdateCounterMetric(log *zerolog.Logger, metricName string, metricValue string) error
	}
)

const (
	defaultBadRequestMessage = "bad request"
)

type Controller struct {
	metricsService MetricsService
}

func New(metricsService MetricsService) *Controller {
	return &Controller{
		metricsService: metricsService,
	}
}

func (c *Controller) Update(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	switch ctx.Param("metricType") {
	case models.GaugeMetricType:
		err := c.metricsService.UpdateGaugeMetric(&log, ctx.Param("metricName"), ctx.Param("metricValue"))
		if err != nil {
			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}
	case models.CounterMetricType:
		err := c.metricsService.UpdateCounterMetric(&log, ctx.Param("metricName"), ctx.Param("metricValue"))
		if err != nil {
			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}
	default:
		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	return ctx.String(http.StatusOK, "metric was updated successfully")
}

func (c *Controller) GetMetric(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	metricType := ctx.Param("metricType")
	if metricType != models.GaugeMetricType && metricType != models.CounterMetricType {
		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	metricValue := c.metricsService.GetMetric(metricType, ctx.Param("metricName"))
	if metricValue == nil {
		return ctx.String(http.StatusNotFound, "not found metric")
	}

	return ctx.JSON(http.StatusOK, *metricValue)
}

func (c *Controller) GetAllMetrics(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	return ctx.JSON(http.StatusOK, c.metricsService.GetAllMetrics())
}
