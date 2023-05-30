package controllers

import (
	"compress/gzip"
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
		GetGaugeMetric(metricName string) *float64
		GetCounterMetric(metricName string) *int64
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

		return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
	}

	switch input.MType {
	case models.GaugeMetricType:
		if nil == input.Value {
			log.Warn().Msg("value of gauge metric is empty")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

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
		if nil == input.Delta {
			log.Warn().Msg("value of counter metric is empty")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

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

		return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
	}

	switch input.MType {
	case models.GaugeMetricType:
		input.Value = c.metricsService.GetGaugeMetric(input.ID)
	case models.CounterMetricType:
		input.Delta = c.metricsService.GetCounterMetric(input.ID)
	default:
		log.Warn().Str("metric_type", input.MType).Msg("invalid metric type")

		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	if input.Delta == nil && input.Value == nil {
		return ctx.String(http.StatusNotFound, "not found metric")
	}

	return ctx.JSON(http.StatusOK, input)
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

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
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

	var metricDelta *float64
	var metricValue *int64

	switch ctx.Param("metricType") {
	case models.GaugeMetricType:
		metricDelta = c.metricsService.GetGaugeMetric(ctx.Param("metricName"))
		if metricDelta == nil {
			return ctx.String(http.StatusNotFound, "not found metric")
		}

		return ctx.JSON(http.StatusOK, *metricDelta)
	case models.CounterMetricType:
		metricValue = c.metricsService.GetCounterMetric(ctx.Param("metricName"))
		if metricValue == nil {
			return ctx.String(http.StatusNotFound, "not found metric")
		}

		return ctx.JSON(http.StatusOK, *metricValue)
	}

	log.Warn().Str("metric_type", ctx.Param("metricType")).Msg("invalid metric type")

	return ctx.String(http.StatusBadRequest, "invalid type of metric")
}

func (c *Controller) GetAllMetrics(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.metricsService.GetAllMetrics())
}

func unmarshalRequestBody(ctx echo.Context) (*MetricInput, error) {
	var reader io.Reader

	if ctx.Request().Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(ctx.Request().Body)
		if err != nil {
			http.Error(ctx.Response().Writer, err.Error(), http.StatusInternalServerError)

			return nil, err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = ctx.Request().Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(ctx.Response().Writer, err.Error(), http.StatusInternalServerError)

		return nil, err
	}

	input := &MetricInput{}
	err = json.Unmarshal(body, input)
	if err != nil {
		return nil, err
	}

	return input, err
}
