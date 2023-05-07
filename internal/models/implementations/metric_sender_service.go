package implementations

import (
	"encoding/json"
	"fmt"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services"
	"io"
	"net/http"
	"strings"

	"github.com/lenarsaitov/metrics-tpl/internal/models/services/metricsender"
)

type MetricSenderModel struct {
	pollCount           int64
	remoteServerAddress string
}

var _ metricsender.Service = &MetricSenderModel{}

func NewMetricSenderModel(
	remoteServerAddress string,
) *MetricSenderModel {
	if !strings.Contains(remoteServerAddress, `://`) {
		remoteServerAddress = "http://" + remoteServerAddress
	}

	return &MetricSenderModel{
		remoteServerAddress: remoteServerAddress,
	}
}

func (m *MetricSenderModel) SendReplaceGauge(name string, value float64) error {
	return m.send(fmt.Sprintf("/update/%s/%s/%f", services.GaugeMetricType, name, value))
}

func (m *MetricSenderModel) SendAddCounter(name string, value int64) error {
	return m.send(fmt.Sprintf("/update/%s/%s/%d", services.CounterMetricType, name, value))
}

func (m *MetricSenderModel) send(urlPath string) error {
	request, err := http.NewRequest(http.MethodPost, m.remoteServerAddress+urlPath, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := ServerResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccess response status: %d, text: %s", resp.StatusCode, response.Response.Text)
	}

	return nil
}
