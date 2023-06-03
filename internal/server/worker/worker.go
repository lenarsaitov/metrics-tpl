package worker

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

type (
	MetricsService interface {
		GetAll(ctx context.Context) (models.Metrics, error)
		UpdateGaugeMetric(ctx context.Context, metricName string, gaugeValue float64) error
		UpdateCounterMetric(ctx context.Context, metricName string, counterValue int64) (int64, error)
	}
)

type Worker struct {
	metricsService  MetricsService
	fileStoragePath string
	storeInterval   time.Duration
}

func New(
	metricsService MetricsService,
	storeInterval int,
	fileStoragePath string,
) *Worker {
	return &Worker{
		metricsService:  metricsService,
		storeInterval:   time.Second * time.Duration(storeInterval),
		fileStoragePath: fileStoragePath,
	}
}

func (r *Worker) Run(ctx context.Context, restore bool) {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()

	if len(r.fileStoragePath) == 0 {
		log.Info().Msg("there is empty file storage path, stop worker")

		return
	}

	if restore {
		log.Info().Msg("there is restore == true, load saved data from file to server")

		if err := r.restoreMetrics(); err != nil {
			log.Error().Err(err).Msg("failed to restore metrics")

			return
		}
	}

	if r.storeInterval == 0 {
		r.savingMetrics(&log)

		return
	}

	go r.savingMetricsAsync(ctx, &log)
}

func (r *Worker) restoreMetrics() error {
	file, err := os.OpenFile(r.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	metrics := &models.Metrics{}

	err = json.Unmarshal(data, metrics)
	if err != nil {
		return err
	}

	for _, metric := range metrics.GaugeMetrics {
		err = r.metricsService.UpdateGaugeMetric(context.Background(), metric.Name, metric.Value)
		if err != nil {
			return err
		}
	}

	for _, metric := range metrics.CounterMetrics {
		_, err = r.metricsService.UpdateCounterMetric(context.Background(), metric.Name, metric.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Worker) savingMetricsAsync(ctx context.Context, logger *zerolog.Logger) {
	log := logger

	ticker := time.NewTicker(r.storeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics, err := r.metricsService.GetAll(context.Background())
			if err != nil {
				log.Error().Err(err).Msg("failed to get all metrics")

				return
			}

			data, err := json.Marshal(metrics)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal metrics data")

				return
			}

			file, err := os.OpenFile(r.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Error().Err(err).Msg("failed to open file")

				return
			}

			_, err = file.Write(data)
			if err != nil {
				log.Error().Err(err).Msg("failed to write data to file")

				return
			}

			log.Info().Str("file_path", r.fileStoragePath).Msg("metrics data saved to file")
		}
	}
}

func (r *Worker) savingMetrics(logger *zerolog.Logger) {
	log := logger

	metrics, err := r.metricsService.GetAll(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to get all metrics")

		return
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal metrics data")

		return
	}

	file, err := os.OpenFile(r.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error().Err(err).Msg("failed to open file")

		return
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("failed to write data to file")

		return
	}

	log.Info().Str("file_path", r.fileStoragePath).Msg("metrics data saved to file")
}
