package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"sync"
	"time"
)

type PollStorage interface {
	GetPoll() models.Metrics
}

type MetricOutput struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MetricsService struct {
	client              http.Client
	metricPollService   PollStorage
	polledMetrics       models.Metrics
	remoteServerAddress string
	pollCount           int64
	pollInterval        time.Duration
	reportInterval      time.Duration
	mu                  *sync.Mutex
}

func NewMetricsService(
	metricPollService PollStorage,
	remoteServerAddress string,
	pollInterval int,
	reportInterval int,
) *MetricsService {
	client := http.Client{}

	return &MetricsService{
		client:              client,
		metricPollService:   metricPollService,
		remoteServerAddress: remoteServerAddress,
		polledMetrics:       models.Metrics{GaugeMetrics: make([]models.GaugeMetric, 0), CounterMetrics: make([]models.CounterMetric, 0)},
		pollInterval:        time.Duration(pollInterval) * time.Second,
		reportInterval:      time.Duration(reportInterval) * time.Second,
		mu:                  &sync.Mutex{},
	}
}

func (s *MetricsService) PollAndReport(ctx context.Context, log *zerolog.Logger) {
	log.Info().Msg("start poll metrics...")
	go s.pollMetrics(ctx, log)

	log.Info().Msg("start report metrics...")
	s.reportMetrics(ctx, log)
}

func (s *MetricsService) pollMetrics(ctx context.Context, log *zerolog.Logger) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("poll ticker stopped by ctx")

			return
		case <-ticker.C:
			metrics := s.metricPollService.GetPoll()
			s.mu.Lock()
			s.polledMetrics = metrics
			s.mu.Unlock()
		}
	}
}

func (s *MetricsService) reportMetrics(ctx context.Context, log *zerolog.Logger) {
	ticker := time.NewTicker(s.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("report ticker stopped by ctx")

			return
		case <-ticker.C:
			s.mu.Lock()
			gaugeMetrics := s.polledMetrics.GaugeMetrics
			counterMetrics := s.polledMetrics.CounterMetrics
			s.mu.Unlock()

			for _, metric := range gaugeMetrics {
				input := &MetricOutput{ID: metric.Name, Value: &metric.Value, MType: models.GaugeMetricType}
				body, err := json.Marshal(input)
				if err != nil {
					log.Error().Err(err).Msg("failed to marshal request body")

					return
				}

				err = s.send(ctx, body)
				if err != nil {
					log.Error().Err(err).RawJSON("body", body).Msg("failed to report gauge metric")

					return
				}
			}
			for _, metric := range counterMetrics {
				input := &MetricOutput{ID: metric.Name, Delta: &metric.Value, MType: models.CounterMetricType}
				body, err := json.Marshal(input)
				if err != nil {
					log.Error().Err(err).RawJSON("body", body).Msg("failed to marshal request body")

					return
				}

				err = s.send(ctx, body)
				if err != nil {
					log.Error().Err(err).RawJSON("body", body).Msg("failed to report counter metric")

					return
				}
			}

			log.Info().Interface("metrics", s.polledMetrics).Msg("reported metrics")
		}
	}
}

func (s *MetricsService) send(ctx context.Context, body []byte) error {
	reader := bytes.NewReader(body)

	//var buf bytes.Buffer
	//
	//g := gzip.NewWriter(&buf)
	//if _, err := g.Write(body); err != nil {
	//	return err
	//}
	//
	//if err := g.Close(); err != nil {
	//	return err
	//}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.remoteServerAddress+"/update/", reader)
	if err != nil {
		return err
	}

	request.Close = true
	request.Header.Set("Content-type", "application/json")
	//request.Header.Set("Content-Encoding", "gzip")

	resp, err := s.client.Do(request)
	if err != nil {
		if errors.As(err, &io.EOF) {
			return nil
		}

		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccess response, url: %s, status: %d", request.URL.String(), resp.StatusCode)
	}

	return nil
}
