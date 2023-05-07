package server

import (
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server/update"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/memstorage"
)

type Controller struct {
	updateHandler *update.Handler
}

func NewController(memStorageService memstorage.Service) *Controller {
	return &Controller{
		updateHandler: update.NewHandler(memStorageService),
	}
}
