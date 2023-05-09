package controllers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	logger "github.com/rs/zerolog/log"
)

type Controller struct {
	metricsUseCase MetricsServerUseCase
}

func New(metricsUseCase MetricsServerUseCase) *Controller {
	return &Controller{
		metricsUseCase: metricsUseCase,
	}
}

func (c *Controller) Update(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()
	rsp := NewResponder(ctx)

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	switch ctx.Param("metricType") {
	case models.GaugeMetricType:
		err := c.metricsUseCase.UpdateGaugeMetric(&log, ctx.Param("metricName"), ctx.Param("metricValue"))
		if err != nil {
			return rsp.BadRequest(defaultBadRequestMessage)
		}
	case models.CounterMetricType:
		err := c.metricsUseCase.UpdateCounterMetric(&log, ctx.Param("metricName"), ctx.Param("metricValue"))
		if err != nil {
			return rsp.BadRequest(defaultBadRequestMessage)
		}
	default:
		return rsp.BadRequest("invalid type of metric")
	}

	return rsp.OK("metric was updated successfully")
}

func (c *Controller) GetMetric(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()
	rsp := NewResponder(ctx)

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	metricType := ctx.Param("metricType")
	if metricType != models.GaugeMetricType && metricType != models.CounterMetricType {
		return rsp.BadRequest("invalid type of metric")
	}

	metricValue := c.metricsUseCase.GetMetric(metricType, ctx.Param("metricName"))
	if metricValue == nil {
		return rsp.NotFound("not found metric")
	}

	return rsp.OKWithBody(*metricValue)
}

func (c *Controller) GetMetrics(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()
	rsp := NewResponder(ctx)

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	return rsp.OKWithBody(c.metricsUseCase.GetAllMetrics())
}
