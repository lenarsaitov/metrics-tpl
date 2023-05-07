package responder

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type Responder struct {
	log *zerolog.Logger
	w   http.ResponseWriter
}

func NewResponder(log *zerolog.Logger, w http.ResponseWriter) *Responder {
	return &Responder{
		log: log,
		w:   w,
	}
}

func (r *Responder) OK(responseMessage string) {
	r.send(http.StatusOK, fmt.Sprintf(requestTemplate, responseMessage))
}

func (r *Responder) BadRequest(responseMessage string) {
	r.send(http.StatusBadRequest, fmt.Sprintf(requestTemplate, responseMessage))
}

func (r *Responder) InternalError() {
	r.send(http.StatusInternalServerError, fmt.Sprintf(requestTemplate, defaultInternalErrorResponseMessage))
}

func (r *Responder) send(status int, responseMessage string) {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(status)

	_, err := r.w.Write([]byte(responseMessage))
	if err != nil {
		r.log.Error().Err(err).Msg("write json to response")
		r.w.WriteHeader(http.StatusInternalServerError)
	}
}
