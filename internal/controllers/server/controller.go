package server

import (
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server/getmetric"
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server/getmetrics"
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server/update"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/memstorage"
)

type Controller struct {
	updateHandler     *update.Handler
	getMetricHandler  *getmetric.Handler
	getMetricsHandler *getmetrics.Handler
}

func NewController(memStorageService memstorage.Service) *Controller {
	return &Controller{
		updateHandler:     update.NewHandler(memStorageService),
		getMetricHandler:  getmetric.NewHandler(memStorageService),
		getMetricsHandler: getmetrics.NewHandler(memStorageService),
	}
}
