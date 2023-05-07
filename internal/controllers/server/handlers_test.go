package server

import (
	"github.com/labstack/echo"
	"github.com/lenarsaitov/metrics-tpl/internal/models/implementations"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdate(t *testing.T) {
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
			}{metricType: services.GaugeMetricType, metricName: "Alloc", metricValue: "123", method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
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

			serverController := NewController(implementations.NewMemStorageModel())
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.request.method, "/update/:metricType/:metricName/:metricValue", nil)

			ctx := e.NewContext(request, w)
			ctx.SetParamNames("metricType", "metricName", "metricValue")
			ctx.SetParamValues(test.request.metricType, test.request.metricName, test.request.metricValue)

			err := serverController.Update(ctx)
			require.Nil(t, err)

			response := w.Result()
			defer response.Body.Close()

			require.Equal(t, response.StatusCode, test.want.statusCode, "incorrect status")
		})
	}
}

func TestGetMetric(t *testing.T) {
	var tests = []struct {
		name    string
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
			name: "test success case",
			request: struct {
				metricType string
				metricName string
				method     string
			}{metricType: services.GaugeMetricType, metricName: "Alloc", method: http.MethodGet},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
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

			memstorageModel := implementations.NewMemStorageModel()
			memstorageModel.ReplaceGauge(test.request.metricName, rand.Float64())

			serverController := NewController(memstorageModel)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.request.method, "/update/:metricType/:metricName", nil)

			ctx := e.NewContext(request, w)
			ctx.SetParamNames("metricType", "metricName")
			ctx.SetParamValues(test.request.metricType, test.request.metricName)

			err := serverController.GetMetric(ctx)
			require.Nil(t, err)

			response := w.Result()
			defer response.Body.Close()

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

			memstorageModel := implementations.NewMemStorageModel()
			memstorageModel.ReplaceGauge(services.GaugeMetricType, rand.Float64())

			serverController := NewController(memstorageModel)
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.requestMethod, "/", nil)

			ctx := e.NewContext(request, w)
			err := serverController.GetMetrics(ctx)
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
