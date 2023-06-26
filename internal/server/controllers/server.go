package controllers

import (
	"compress/gzip"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

const (
	defaultInternalErrorMessage = `something going wrong: %s`
	defaultBadRequestMessage    = "bad request"
)

type (
	MetricsService interface {
		GetAll(ctx context.Context) (models.Metrics, error)
		GetGaugeMetric(ctx context.Context, metricName string) (*float64, error)
		GetCounterMetric(ctx context.Context, metricName string) (*int64, error)
		UpdateGaugeMetric(ctx context.Context, metricName string, gaugeValue float64) error
		UpdateCounterMetric(ctx context.Context, metricName string, counterValue int64) (int64, error)
	}
)

type MetricInput struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

type Controller struct {
	metricsService MetricsService
	dataSourceName string
	jwtKey         string
}

func New(metricsService MetricsService, dataSourceName string, jwtKey string) *Controller {
	return &Controller{
		metricsService: metricsService,
		dataSourceName: dataSourceName,
		jwtKey:         jwtKey,
	}
}

func (c *Controller) PingDB(ctx echo.Context) error {
	db, err := sql.Open("pgx", c.dataSourceName)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to postgresql database")

		return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
	}
	defer db.Close()

	err = db.PingContext(context.Background())
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
	}

	return ctx.String(http.StatusOK, "OK")
}

func (c *Controller) Update(ctx echo.Context) error {
	body, err := getRequestBody(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed get body from request")

		return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
	}

	input := &MetricInput{}
	err = json.Unmarshal(body, input)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
	}

	switch input.MType {
	case models.GaugeMetricType:
		if nil == input.Value {
			log.Warn().Msg("value of gauge metric is empty")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

		err = c.metricsService.UpdateGaugeMetric(context.Background(), input.ID, *input.Value)
		if err != nil {
			log.Error().Err(err).Msg("failed to update gauge metric")

			return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
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

		value, err := c.metricsService.UpdateCounterMetric(context.Background(), input.ID, *input.Delta)
		if err != nil {
			log.Error().Err(err).Msg("failed to update counter metric")

			return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
		}

		actualCounterValue := float64(value)
		input.Value = &actualCounterValue

		log.Info().
			Str("metric_name", input.ID).
			Int64("counter_value", *input.Delta).
			Msg("counter was added successfully")
	default:
		log.Warn().Msg("unknow metric type")

		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	return ctx.JSON(http.StatusOK, *input)
}

func (c *Controller) Updates(ctx echo.Context) error {
	body, err := getRequestBody(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed get body from request")

		return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
	}

	if c.jwtKey != "" {
		h := sha256.New()
		_, err = h.Write(append(body, []byte(c.jwtKey)...))
		if err != nil {
			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

		if base64.StdEncoding.EncodeToString(h.Sum(nil)) != ctx.Request().Header.Get("HashSHA256") {
			log.Error().Msg("hash of body is not valid")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}
	}

	var inputMetrics []MetricInput
	err = json.Unmarshal(body, &inputMetrics)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
	}

	if len(inputMetrics) == 0 {
		return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
	}

	for _, input := range inputMetrics {
		switch input.MType {
		case models.GaugeMetricType:
			if nil == input.Value {
				log.Warn().Msg("value of gauge metric is empty")

				return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
			}

			err = c.metricsService.UpdateGaugeMetric(context.Background(), input.ID, *input.Value)
			if err != nil {
				log.Error().Err(err).Msg("failed to update gauge metric")

				return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
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

			value, err := c.metricsService.UpdateCounterMetric(context.Background(), input.ID, *input.Delta)
			if err != nil {
				log.Error().Err(err).Msg("failed to update counter metric")

				return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
			}

			actualCounterValue := float64(value)
			input.Value = &actualCounterValue

			log.Info().
				Str("metric_name", input.ID).
				Int64("counter_value", *input.Delta).
				Msg("counter was added successfully")
		default:
			log.Warn().Msg("unknow metric type")

			return ctx.String(http.StatusBadRequest, "invalid type of metric")
		}
	}

	return ctx.JSON(http.StatusOK, inputMetrics[0])
}

func (c *Controller) GetMetric(ctx echo.Context) error {
	body, err := getRequestBody(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed get body from request")

		return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
	}

	input := &MetricInput{}
	err = json.Unmarshal(body, input)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
	}

	switch input.MType {
	case models.GaugeMetricType:
		value, err := c.metricsService.GetGaugeMetric(context.Background(), input.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get gauge metric")

			return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
		}

		input.Value = value
	case models.CounterMetricType:
		delta, err := c.metricsService.GetCounterMetric(context.Background(), input.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get counter metric")

			return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
		}

		input.Delta = delta
	default:
		log.Warn().Str("metric_type", input.MType).Msg("invalid type of metric")

		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	if input.Delta == nil && input.Value == nil {
		log.Warn().Str("metric_name", input.ID).Str("metric_type", input.MType).Msg("not found metric")

		return ctx.String(http.StatusNotFound, "not found metric")
	}

	return ctx.JSON(http.StatusOK, *input)
}

func (c *Controller) UpdatePath(ctx echo.Context) error {
	switch ctx.Param("metricType") {
	case models.GaugeMetricType:
		gaugeValue, err := strconv.ParseFloat(ctx.Param("metricValue"), 64)
		if err != nil {
			log.Error().Err(err).Msg("invalid metric value")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

		err = c.metricsService.UpdateGaugeMetric(context.Background(), ctx.Param("metricName"), gaugeValue)
		if err != nil {
			log.Error().Err(err).Msg("failed to update gauge metric")

			return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
		}

		log.Info().
			Str("metric_name", ctx.Param("metricName")).
			Float64("gauge_value", gaugeValue).
			Msg("gauge was replaced successfully")
	case models.CounterMetricType:
		countValue, err := strconv.Atoi(ctx.Param("metricValue"))
		if err != nil {
			log.Error().Err(err).Msg("invalid metric value")

			return ctx.String(http.StatusBadRequest, defaultBadRequestMessage)
		}

		_, err = c.metricsService.UpdateCounterMetric(context.Background(), ctx.Param("metricName"), int64(countValue))
		if err != nil {
			log.Error().Err(err).Msg("failed to update counter metric")

			return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
		}
		log.Info().
			Str("metric_name", ctx.Param("metricName")).
			Int("counter_value", countValue).
			Msg("counter was added successfully")
	default:
		return ctx.String(http.StatusBadRequest, "invalid type of metric")
	}

	return ctx.String(http.StatusOK, "metric was updated successfully")
}

func (c *Controller) GetMetricPath(ctx echo.Context) error {
	var metricDelta *float64
	var metricValue *int64
	var err error

	switch ctx.Param("metricType") {
	case models.GaugeMetricType:
		metricDelta, err = c.metricsService.GetGaugeMetric(context.Background(), ctx.Param("metricName"))
		if err != nil {
			log.Error().Err(err).Msg("failed to get gauge metric")

			return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
		}

		if metricDelta == nil {
			return ctx.String(http.StatusNotFound, "not found metric")
		}

		return ctx.JSON(http.StatusOK, *metricDelta)
	case models.CounterMetricType:
		metricValue, err = c.metricsService.GetCounterMetric(context.Background(), ctx.Param("metricName"))
		if err != nil {
			log.Error().Err(err).Msg("failed to get counter metric")

			return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
		}

		if metricValue == nil {
			return ctx.String(http.StatusNotFound, "not found metric")
		}

		return ctx.JSON(http.StatusOK, *metricValue)
	}

	log.Warn().Str("metric_type", ctx.Param("metricType")).Msg("invalid metric type")

	return ctx.String(http.StatusBadRequest, "invalid type of metric")
}

func (c *Controller) GetAllMetrics(ctx echo.Context) error {
	metrics, err := c.metricsService.GetAll(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to get metrics")

		return ctx.String(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		return ctx.HTML(http.StatusInternalServerError, fmt.Sprintf(defaultInternalErrorMessage, err.Error()))
	}

	return ctx.HTML(http.StatusOK, string(data))
}

func getRequestBody(ctx echo.Context) ([]byte, error) {
	var reader io.Reader

	if ctx.Request().Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(ctx.Request().Body)
		if err != nil {
			return nil, err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = ctx.Request().Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return body, err
}
