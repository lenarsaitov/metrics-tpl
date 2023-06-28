package services

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lenarsaitov/metrics-tpl/internal/agent/models"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"sync"
	"syscall"
	"time"
)

type PollStorage interface {
	PutCommonPoll()
	PutPsutilPoll()

	GetPoll() models.Metrics
}

type MetricOutput struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

type MetricsService struct {
	jobs                chan models.Metrics
	metricPollService   PollStorage
	mu                  *sync.Mutex
	client              http.Client
	remoteServerAddress string
	jwtKey              string
	pollCount           int64
	pollInterval        time.Duration
	reportInterval      time.Duration
	connectionAttempt   int
}

func NewMetricsService(
	jobs chan models.Metrics,
	metricPollService PollStorage,
	remoteServerAddress string,
	pollInterval int,
	reportInterval int,
	jwtKey string,
) *MetricsService {
	client := http.Client{}

	return &MetricsService{
		jobs:                jobs,
		client:              client,
		metricPollService:   metricPollService,
		remoteServerAddress: remoteServerAddress,
		pollInterval:        time.Duration(pollInterval) * time.Second,
		reportInterval:      time.Duration(reportInterval) * time.Second,
		mu:                  &sync.Mutex{},
		jwtKey:              jwtKey,
	}
}

func (s *MetricsService) PutCommonPollWorker(ctx context.Context, log *zerolog.Logger) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("poll ticker stopped by ctx")

			return
		case <-ticker.C:
			s.metricPollService.PutCommonPoll()
		}
	}
}

func (s *MetricsService) PutPsutilPollWorker(ctx context.Context, log *zerolog.Logger) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("poll ticker stopped by ctx")

			return
		case <-ticker.C:
			s.metricPollService.PutPsutilPoll()
		}
	}
}

func (s *MetricsService) WriteToChanWorker(ctx context.Context, log *zerolog.Logger) {
	ticker := time.NewTicker(s.reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("write to chan worker ticker stopped by ctx")

			return
		case <-ticker.C:
			s.jobs <- s.metricPollService.GetPoll()
		}
	}
}

func (s *MetricsService) SendWorker(ctx context.Context, log *zerolog.Logger) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("report ticker stopped by ctx")

			return
		case metrics := <-s.jobs:
			var input []MetricOutput

			for _, metric := range metrics.GaugeMetrics {
				input = append(input, MetricOutput{ID: metric.Name, Value: &metric.Value, MType: models.GaugeMetricType})
			}

			for _, metric := range metrics.CounterMetrics {
				input = append(input, MetricOutput{ID: metric.Name, Delta: &metric.Value, MType: models.CounterMetricType})
			}

			body, err := json.Marshal(input)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal request body")

				return
			}

			err = s.report(ctx, log, body)
			if err != nil {
				log.Error().Err(err).Msg("failed to report metrics")

				return
			}

			log.Info().
				Interface("gauge_metrics", metrics.GaugeMetrics).
				Interface("counter_metrics", metrics.CounterMetrics).
				Msg("reported metrics")
		}
	}
}

func (s *MetricsService) report(ctx context.Context, log *zerolog.Logger, body []byte) error {
	for {
		err := s.send(ctx, body)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) {
				log.Warn().Int("next_attempt_after_seconds", 2*s.connectionAttempt+1).Msg("couldn't send request, try next attempt")
				time.Sleep(time.Duration(2*s.connectionAttempt+1) * time.Second)

				s.connectionAttempt++

				continue
			}
			log.Error().Err(err).RawJSON("request_body", body).Msg("failed to report gauge metric")

			return err
		}

		return nil
	}
}

func (s *MetricsService) send(ctx context.Context, reqBody []byte) error {
	var buf bytes.Buffer

	g := gzip.NewWriter(&buf)
	if _, err := g.Write(reqBody); err != nil {
		return err
	}

	if err := g.Close(); err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.remoteServerAddress+"/updates/", &buf)
	if err != nil {
		return err
	}

	request.Close = true
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")

	if s.jwtKey != "" {
		h := sha256.New()
		_, err = h.Write(append(reqBody, []byte(s.jwtKey)...))
		if err != nil {
			return err
		}

		request.Header.Set("HashSHA256", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	}

	resp, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccess response, url: %s, body: %s, status: %d", request.URL.String(), string(respBody), resp.StatusCode)
	}

	return nil
}
