package server

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/lenarsaitov/metrics-tpl/internal/controllers/server/update"
	"github.com/lenarsaitov/metrics-tpl/internal/models/services/memstorage"
	"github.com/lenarsaitov/metrics-tpl/internal/responder"
	logger "github.com/rs/zerolog/log"
)

type Controller struct {
	updateHandler *update.Handler
}

func NewController(memStorageService memstorage.Service) *Controller {
	return &Controller{
		updateHandler: update.NewHandler(memStorageService),
	}
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	log := logger.With().Str("request_id", uuid.New().String()).Logger()
	rsp := responder.NewResponder(&log, w)

	if r.Method != http.MethodPost {
		rsp.BadRequest("bad method of request, need post")

		return
	}

	paths := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(paths) < 2 || paths[0] != "update" || (paths[1] != memstorage.GaugeMetricType && paths[1] != memstorage.CounterMetricType) {
		rsp.BadRequest("invalid url path structure")

		return
	}

	if len(paths) != 4 {
		rsp.NotFound("not found metric, doesnt have name of it")

		return
	}

	log.Info().Str("url", r.URL.String()).Msg("url of request")

	c.updateHandler.Handle(
		&log,
		rsp,
		&update.Input{
			MetricType:  paths[1],
			MetricName:  paths[2],
			MetricValue: paths[3],
		},
	)
}
