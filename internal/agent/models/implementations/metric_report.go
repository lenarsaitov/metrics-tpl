package implementations

import (
	"encoding/json"
	"fmt"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"io"
	"net/http"
	"strings"
)

type MetricReportModel struct {
	pollCount           int64
	remoteServerAddress string
}

var _ models.MetricReport = &MetricReportModel{}

func NewMetricReportModel(
	remoteServerAddress string,
) *MetricReportModel {
	if !strings.Contains(remoteServerAddress, `://`) {
		remoteServerAddress = "http://" + remoteServerAddress
	}

	return &MetricReportModel{
		remoteServerAddress: remoteServerAddress,
	}
}

func (m *MetricReportModel) ReportGaugeMetric(name string, value float64) error {
	return m.send(fmt.Sprintf("/update/%s/%s/%f", models.GaugeMetricType, name, value))
}

func (m *MetricReportModel) ReportCounterMetric(name string, value int64) error {
	return m.send(fmt.Sprintf("/update/%s/%s/%d", models.CounterMetricType, name, value))
}

func (m *MetricReportModel) send(urlPath string) error {
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
