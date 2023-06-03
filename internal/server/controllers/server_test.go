package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/lenarsaitov/metrics-tpl/internal/server/repository/inmemory"
	"github.com/lenarsaitov/metrics-tpl/internal/server/services"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdatePath(t *testing.T) {
	var tests = []struct {
		name    string
		request struct {
			metricType  string
			metricName  string
			metricValue string
			method      string
		}
		want struct {
			statusCode int
		}
	}{
		{
			name: "test success case",
			request: struct {
				metricType  string
				metricName  string
				metricValue string
				method      string
			}{metricType: models.GaugeMetricType, metricName: "Alloc", metricValue: "123", method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
		},
		{
			name: "test negative case, invalid metric value",
			request: struct {
				metricType  string
				metricName  string
				metricValue string
				method      string
			}{metricType: models.GaugeMetricType, metricName: "Alloc", metricValue: "123ssss", method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusBadRequest},
		},
		{
			name: "test negative case, incorrect metric type",
			request: struct {
				metricType  string
				metricName  string
				metricValue string
				method      string
			}{metricType: "counter111", metricName: "Alloc", metricValue: "123", method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusBadRequest},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			useMetrics := services.NewMetricsService(inmemory.NewPollStorage())
			serverController := New("", useMetrics)

			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.request.method, "/update/:metricType/:metricName/:metricValue", nil)

			ctx := e.NewContext(request, w)
			ctx.SetParamNames("metricType", "metricName", "metricValue")
			ctx.SetParamValues(test.request.metricType, test.request.metricName, test.request.metricValue)

			err := serverController.UpdatePath(ctx)
			require.Nil(t, err)

			response := w.Result()
			defer response.Body.Close()
			require.Equal(t, response.StatusCode, test.want.statusCode, "incorrect status")
		})
	}
}
func TestGetMetricPath(t *testing.T) {
	var tests = []struct {
		name           string
		preparedMetric *struct {
			metricType string
			metricName string
		}
		request struct {
			metricType string
			metricName string
			method     string
		}
		want struct {
			statusCode int
		}
	}{
		{
			name: "test success case, gauge",
			preparedMetric: &struct {
				metricType string
				metricName string
			}{metricType: models.GaugeMetricType, metricName: "Alloc"},
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: models.GaugeMetricType, metricName: "Alloc", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
		},
		{
			name: "test success case, counter",
			preparedMetric: &struct {
				metricType string
				metricName string
			}{metricType: models.CounterMetricType, metricName: "Counter"},
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: models.CounterMetricType, metricName: "Counter", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
		},
		{
			name: "test negative case, not found",
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: models.CounterMetricType, metricName: "Alloc", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusNotFound},
		},
		{
			name: "test negative case, incorrect metric type",
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: "counter111", metricName: "Alloc", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusBadRequest},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			memStorageModel := inmemory.NewPollStorage()
			if test.preparedMetric != nil {
				switch test.preparedMetric.metricType {
				case models.GaugeMetricType:
					err := memStorageModel.ReplaceGauge(context.Background(), test.request.metricName, rand.Float64())
					require.Nil(t, err)
				case models.CounterMetricType:
					_, err := memStorageModel.AddCounter(context.Background(), test.request.metricName, rand.Int63())
					require.Nil(t, err)
				}
			}

			useMetrics := services.NewMetricsService(memStorageModel)

			serverController := New("", useMetrics)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.request.method, "/value/:metricType/:metricName", nil)

			ctx := e.NewContext(request, w)
			ctx.SetParamNames("metricType", "metricName")
			ctx.SetParamValues(test.request.metricType, test.request.metricName)

			err := serverController.GetMetricPath(ctx)
			require.Nil(t, err)

			response := w.Result()
			defer response.Body.Close()
			res, err := io.ReadAll(response.Body)
			require.Nil(t, err)
			require.NotEmpty(t, res)
			require.Equal(t, response.StatusCode, test.want.statusCode, "incorrect status")
		})
	}
}

func TestUpdateGauge(t *testing.T) {
	var tests = []struct {
		name    string
		request struct {
			metricType  string
			metricName  string
			method      string
			metricValue float64
		}
		want struct {
			statusCode int
		}
	}{
		{
			name: "test success case",
			request: struct {
				metricType  string
				metricName  string
				method      string
				metricValue float64
			}{metricType: models.GaugeMetricType, metricName: "Alloc", metricValue: 123, method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
		},
		{
			name: "test negative case, incorrect metric type",
			request: struct {
				metricType  string
				metricName  string
				method      string
				metricValue float64
			}{metricType: "counter111", metricName: "Alloc", metricValue: 333, method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusBadRequest},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()

			useMetrics := services.NewMetricsService(inmemory.NewPollStorage())
			serverController := New("", useMetrics)

			input := &MetricInput{ID: test.request.metricName, MType: test.request.metricType, Value: &test.request.metricValue}
			body, err := json.Marshal(input)
			require.Nil(t, err)

			reader := bytes.NewReader(body)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.request.method, "/update", reader)
			request.Header.Set("Content-type", "application/json")

			ctx := e.NewContext(request, w)

			err = serverController.Update(ctx)
			require.Nil(t, err)

			response := w.Result()
			defer response.Body.Close()

			require.Equal(t, response.StatusCode, test.want.statusCode, "incorrect status")
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	var tests = []struct {
		name    string
		request struct {
			metricType  string
			metricName  string
			method      string
			metricValue int64
		}
		want struct {
			statusCode int
		}
	}{
		{
			name: "test success case",
			request: struct {
				metricType  string
				metricName  string
				method      string
				metricValue int64
			}{metricType: models.CounterMetricType, metricName: "Alloc", metricValue: 123, method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
		},
		{
			name: "test negative case, incorrect metric type",
			request: struct {
				metricType  string
				metricName  string
				method      string
				metricValue int64
			}{metricType: "counter111", metricName: "Counter", metricValue: 333, method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusBadRequest},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()

			useMetrics := services.NewMetricsService(inmemory.NewPollStorage())
			serverController := New("", useMetrics)

			input := &MetricInput{ID: test.request.metricName, MType: test.request.metricType, Delta: &test.request.metricValue}
			body, err := json.Marshal(input)
			require.Nil(t, err)

			reader := bytes.NewReader(body)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.request.method, "/update", reader)
			request.Header.Set("Content-type", "application/json")

			ctx := e.NewContext(request, w)

			err = serverController.Update(ctx)
			require.Nil(t, err)

			response := w.Result()
			defer response.Body.Close()

			require.Equal(t, response.StatusCode, test.want.statusCode, "incorrect status")
		})
	}
}

func TestGetMetric(t *testing.T) {
	var tests = []struct {
		name           string
		preparedMetric *struct {
			metricType string
			metricName string
		}
		request struct {
			metricType string
			metricName string
			method     string
		}
		want struct {
			statusCode int
		}
	}{
		{
			name: "test success case, gauge",
			preparedMetric: &struct {
				metricType string
				metricName string
			}{metricType: models.GaugeMetricType, metricName: "Alloc"},
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: models.GaugeMetricType, metricName: "Alloc", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
		},
		{
			name: "test success case, counter",
			preparedMetric: &struct {
				metricType string
				metricName string
			}{metricType: models.CounterMetricType, metricName: "Counter"},
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: models.CounterMetricType, metricName: "Counter", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
		},
		{
			name: "test negative case, not found",
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: models.CounterMetricType, metricName: "Alloc", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusNotFound},
		},
		{
			name: "test negative case, incorrect metric type",
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: "counter111", metricName: "Alloc", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusBadRequest},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()

			memStorageModel := inmemory.NewPollStorage()

			if test.preparedMetric != nil {
				switch test.preparedMetric.metricType {
				case models.GaugeMetricType:
					err := memStorageModel.ReplaceGauge(context.Background(), test.request.metricName, rand.Float64())
					require.Nil(t, err)
				case models.CounterMetricType:
					_, err := memStorageModel.AddCounter(context.Background(), test.request.metricName, rand.Int63())
					require.Nil(t, err)
				}
			}

			useMetrics := services.NewMetricsService(memStorageModel)
			serverController := New("", useMetrics)

			input := &MetricInput{ID: test.request.metricName, MType: test.request.metricType}
			body, err := json.Marshal(input)
			require.Nil(t, err)

			reader := bytes.NewReader(body)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.request.method, "/update", reader)
			request.Header.Set("Content-type", "application/json")

			ctx := e.NewContext(request, w)

			err = serverController.GetMetric(ctx)
			require.Nil(t, err)

			response := w.Result()
			defer response.Body.Close()
			res, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			require.NotEmpty(t, res)
			require.Equal(t, response.StatusCode, test.want.statusCode, "incorrect status")
		})
	}
}

func TestGetAllMetrics(t *testing.T) {
	var tests = []struct {
		name          string
		requestMethod string
		want          struct {
			statusCode int
		}
	}{
		{
			name:          "test success case",
			requestMethod: http.MethodGet,
			want:          struct{ statusCode int }{statusCode: http.StatusOK},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()

			memStorageModel := inmemory.NewPollStorage()
			err := memStorageModel.ReplaceGauge(context.Background(), models.GaugeMetricType, rand.Float64())
			require.Nil(t, err)

			useMetrics := services.NewMetricsService(memStorageModel)

			serverController := New("", useMetrics)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.requestMethod, "/", nil)

			ctx := e.NewContext(request, w)
			err = serverController.GetAllMetrics(ctx)
			require.Nil(t, err)

			response := w.Result()
			defer response.Body.Close()

			res, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			require.NotEmpty(t, res)
			require.Equal(t, response.StatusCode, test.want.statusCode, "incorrect status")
		})
	}
}
