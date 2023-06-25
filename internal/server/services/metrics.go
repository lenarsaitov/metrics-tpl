package services

import (
	"context"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
)

type Storage interface {
	GetAll(ctx context.Context) (models.Metrics, error)

	GetGaugeMetric(ctx context.Context, name string) (*float64, error)
	GetCounterMetric(ctx context.Context, name string) (*int64, error)

	ReplaceGauge(ctx context.Context, name string, value float64) error
	AddCounter(ctx context.Context, name string, value int64) (int64, error)
}

type MetricsService struct {
	storage Storage
}

func NewMetricsService(storageService Storage) *MetricsService {
	return &MetricsService{
		storage: storageService,
	}
}

func (h *MetricsService) GetAll(ctx context.Context) (models.Metrics, error) {
	return h.storage.GetAll(ctx)
}

func (h *MetricsService) GetGaugeMetric(ctx context.Context, metricName string) (*float64, error) {
	return h.storage.GetGaugeMetric(ctx, metricName)
}

func (h *MetricsService) GetCounterMetric(ctx context.Context, metricName string) (*int64, error) {
	return h.storage.GetCounterMetric(ctx, metricName)
}

func (h *MetricsService) UpdateGaugeMetric(ctx context.Context, metricName string, gaugeValue float64) error {
	return h.storage.ReplaceGauge(ctx, metricName, gaugeValue)
}

func (h *MetricsService) UpdateCounterMetric(ctx context.Context, metricName string, counterValue int64) (int64, error) {
	return h.storage.AddCounter(ctx, metricName, counterValue)
}
