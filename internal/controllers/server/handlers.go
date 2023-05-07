package server

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server/getmetric"
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server/update"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services"
	"github.com/lenarsaitov/metrics-tpl/internal/responder"
	logger "github.com/rs/zerolog/log"
)

func (c *Controller) Update(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()
	rsp := responder.NewResponder(&log, ctx)

	metricType := ctx.Param("metricType")
	if metricType != services.GaugeMetricType && metricType != services.CounterMetricType {
		return rsp.BadRequest("invalid type of metric")
	}

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	return c.updateHandler.Handle(
		&log,
		rsp,
		&update.Input{
			MetricType:  metricType,
			MetricName:  ctx.Param("metricName"),
			MetricValue: ctx.Param("metricValue"),
		},
	)
}

func (c *Controller) GetMetric(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()
	rsp := responder.NewResponder(&log, ctx)

	metricType := ctx.Param("metricType")
	if metricType != services.GaugeMetricType && metricType != services.CounterMetricType {
		return rsp.BadRequest("invalid type of metric")
	}

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	return c.getMetricHandler.Handle(
		&log,
		rsp,
		&getmetric.Input{
			MetricType: metricType,
			MetricName: ctx.Param("metricName"),
		},
	)
}

func (c *Controller) GetMetrics(ctx echo.Context) error {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()
	rsp := responder.NewResponder(&log, ctx)

	log.Info().Str("url", ctx.Request().URL.String()).Msg("url of request")

	return c.getMetricsHandler.Handle(
		&log,
		rsp,
	)
}
