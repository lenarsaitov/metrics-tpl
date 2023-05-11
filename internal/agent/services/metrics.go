package services

import (
	"fmt"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"github.com/rs/zerolog"
	"net/http"
	"sync"
	"time"
)

type PollStorage interface {
	GetPoll() models.Metrics
}

type MetricsService struct {
	metricPollService   PollStorage
	polledMetrics       models.Metrics
	remoteServerAddress string
	pollCount           int64
	pollInterval        int
	reportInterval      int
}

func NewMetricsService(
	metricPollService PollStorage,
	remoteServerAddress string,
	pollInterval int,
	reportInterval int,
) *MetricsService {
	return &MetricsService{
		metricPollService:   metricPollService,
		remoteServerAddress: remoteServerAddress,
		polledMetrics:       models.Metrics{GaugeMetrics: make([]models.GaugeMetric, 0), CounterMetrics: make([]models.CounterMetric, 0)},
		pollInterval:        pollInterval,
		reportInterval:      reportInterval,
	}
}

func (s *MetricsService) PollAndReport(log *zerolog.Logger) {
	var mu = &sync.Mutex{}

	log.Info().Msg("start poll metrics...")
	go s.pollMetrics(mu)

	time.Sleep(time.Second * time.Duration(s.pollInterval))

	log.Info().Msg("start report metrics...")
	s.reportMetrics(log)
}

func (s *MetricsService) pollMetrics(mu *sync.Mutex) {
	for {
		metrics := s.metricPollService.GetPoll()

		mu.Lock()
		s.polledMetrics = metrics
		mu.Unlock()

		time.Sleep(time.Second * time.Duration(s.pollInterval))
	}
}

func (s *MetricsService) reportMetrics(log *zerolog.Logger) {
	for {
		for _, metric := range s.polledMetrics.GaugeMetrics {
			err := s.send(fmt.Sprintf("/update/%s/%s/%f", models.GaugeMetricType, metric.Name, metric.Value))
			if err != nil {
				log.Error().Err(err).Msg("failed to report gauge metric")

				return
			}
		}
		for _, metric := range s.polledMetrics.CounterMetrics {
			err := s.send(fmt.Sprintf("/update/%s/%s/%d", models.CounterMetricType, metric.Name, metric.Value))
			if err != nil {
				log.Error().Err(err).Msg("failed to report counter metric")

				return
			}
		}

		log.Info().Interface("metrics", s.polledMetrics).Msg("reported metrics")
		time.Sleep(time.Second * time.Duration(s.reportInterval))
	}
}

func (s *MetricsService) send(urlPath string) error {
	request, err := http.NewRequest(http.MethodPost, s.remoteServerAddress+urlPath, nil)
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
