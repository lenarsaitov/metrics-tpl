package localcache

import (
	"fmt"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"net/http"
)

type MetricReportStorage struct {
	pollCount           int64
	remoteServerAddress string
}

func NewMetricReportStorage(
	remoteServerAddress string,
) *MetricReportStorage {
	return &MetricReportStorage{
		remoteServerAddress: remoteServerAddress,
	}
}

func (m *MetricReportStorage) ReportGaugeMetric(name string, value float64) error {
	return m.send(fmt.Sprintf("/update/%s/%s/%f", models.GaugeMetricType, name, value))
}

func (m *MetricReportStorage) ReportCounterMetric(name string, value int64) error {
	return m.send(fmt.Sprintf("/update/%s/%s/%d", models.CounterMetricType, name, value))
}

func (m *MetricReportStorage) send(urlPath string) error {
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccess response status: %d", resp.StatusCode)
	}

	return nil
}
