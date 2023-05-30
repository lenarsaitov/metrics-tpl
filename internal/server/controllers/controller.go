package controllers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	logger "github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

type (
	MetricsService interface {
		GetAllMetrics() models.Metrics
		GetMetric(metricType, metricName string) *float64
		UpdateGaugeMetric(metricName string, gaugeValue float64) error
		UpdateCounterMetric(metricName string, counterValue int64) error
	}
)

const (
	defaultBadRequestMessage = "bad request"
)

type MetricInput struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Controller struct {
	metricsService MetricsService
}

func New(metricsService MetricsService) *Controller {
	return &Controller{
		metricsService: metricsService,
	}
}

func (c *Controller) Update(ctx echo.Context) error {
	log := logger.With().Logger()

	input, err := unmarshalRequestBody(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal body from request")

		return err
	}

	switch input.MType {
	case models.GaugeMetricType:
		err = c.metricsService.UpdateGaugeMetric(input.ID, *input.Value)
		if err != nil {
			log.Error().Err(err).Msg("failed to update gauge metric")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

		log.Info().
			Str("metric_name", input.ID).
			Float64("gauge_value", *input.Value).
			Msg("gauge was replaced successfully")

	case models.CounterMetricType:
		err = c.metricsService.UpdateCounterMetric(input.ID, *input.Delta)
		if err != nil {
			log.Error().Err(err).Msg("failed to update counter metric")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

		log.Info().
			Str("metric_name", input.ID).
			Int64("counter_value", *input.Delta).
			Msg("counter was added successfully")
	default:
		log.Warn().Msg("unknow metric type")

		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	return ctx.String(http.StatusOK, "metric was updated successfully")
}

func (c *Controller) GetMetric(ctx echo.Context) error {
	log := logger.With().Logger()

	input, err := unmarshalRequestBody(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed unmarshal body from request")

		return err
	}

	if input.MType != models.GaugeMetricType && input.MType != models.CounterMetricType {
		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	metricValue := c.metricsService.GetMetric(input.MType, input.ID)
	if metricValue == nil {
		return ctx.String(http.StatusNotFound, "not found metric")
	}

	return ctx.JSON(http.StatusOK, *metricValue)
}

func (c *Controller) UpdatePath(ctx echo.Context) error {
	log := logger.With().Logger()

	switch ctx.Param("metricType") {
	case models.GaugeMetricType:
		gaugeValue, err := strconv.ParseFloat(ctx.Param("metricValue"), 64)
		if err != nil {
			log.Error().Err(err).Msg("invalid metric value")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

		err = c.metricsService.UpdateGaugeMetric(ctx.Param("metricName"), gaugeValue)
		if err != nil {
			log.Error().Err(err).Msg("failed to update gauge metric")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}
	case models.CounterMetricType:
		countValue, err := strconv.Atoi(ctx.Param("metricValue"))
		if err != nil {
			log.Error().Err(err).Msg("invalid metric value")

			return err
		}

		err = c.metricsService.UpdateCounterMetric(ctx.Param("metricName"), int64(countValue))
		if err != nil {
			log.Error().Err(err).Msg("failed to update counter metric")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}
	default:
		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	return ctx.String(http.StatusOK, "metric was updated successfully")
}

func (c *Controller) GetMetricPath(ctx echo.Context) error {
	log := logger.With().Logger()

	metricType := ctx.Param("metricType")
	if metricType != models.GaugeMetricType && metricType != models.CounterMetricType {
		log.Warn().Str("metric_type", metricType).Msg("invalid metric type")

		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	metricValue := c.metricsService.GetMetric(metricType, ctx.Param("metricName"))
	if metricValue == nil {
		log.Warn().Msg("metric not found")

		return ctx.String(http.StatusNotFound, "not found metric")
	}

	return ctx.JSON(http.StatusOK, *metricValue)
}

func (c *Controller) GetAllMetrics(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.metricsService.GetAllMetrics())
}

func unmarshalRequestBody(ctx echo.Context) (*MetricInput, error) {
	b, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, err
	}

	input := &MetricInput{}
	err = json.Unmarshal(b, input)
	if err != nil {
		return nil, err
	}

	return input, err
}
