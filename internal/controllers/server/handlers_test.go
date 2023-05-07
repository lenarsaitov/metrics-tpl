package server

import (
	"github.com/lenarsaitov/metrics-tpl/internal/models/implementations"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdate(t *testing.T) {
	var tests = []struct {
		name    string
		request struct {
			url    string
			method string
		}
		want struct {
			statusCode int
		}
	}{
		{
			name: "test success case",
			request: struct {
				url    string
				method string
			}{url: "http://localhost/update/gauge/name/1", method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusOK},
		},
		{
			name: "test negative case, incorrect metric type",
			request: struct {
				url    string
				method string
			}{url: "http://localhost/update/gauge1/name/1", method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusBadRequest},
		},
		{
			name: "test negative case, dont have metric name",
			request: struct {
				url    string
				method string
			}{url: "http://localhost/update/gauge/", method: http.MethodPost},
			want: struct{ statusCode int }{statusCode: http.StatusNotFound},
		},
		{
			name: "test negative case, incorrect method of request",
			request: struct {
				url    string
				method string
			}{url: "http://localhost/update/gauge/", method: http.MethodDelete},
			want: struct{ statusCode int }{statusCode: http.StatusBadRequest},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.request.method, test.request.url, nil)

			serverController := NewController(implementations.NewMemStorageModel())
			serverController.Update(w, request)

			response := w.Result()

			require.Equal(t, response.StatusCode, test.want.statusCode, "incorrect status")
		})
	}
}
