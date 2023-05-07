package responder

import (
	"github.com/labstack/echo"
	"net/http"

	"github.com/rs/zerolog"
)

type Responder struct {
	log *zerolog.Logger
	ctx echo.Context
}

func NewResponder(log *zerolog.Logger, ctx echo.Context) *Responder {
	return &Responder{
		log: log,
		ctx: ctx,
	}
}

func (r *Responder) OK(responseMessage string) error {
	return r.send(http.StatusOK, responseMessage)
}

func (r *Responder) OKWithBody(response any) error {
	return r.ctx.JSON(http.StatusOK, response)
}

func (r *Responder) BadRequest(responseMessage string) error {
	return r.send(http.StatusBadRequest, responseMessage)
}

func (r *Responder) NotFound(responseMessage string) error {
	return r.send(http.StatusNotFound, responseMessage)
}

func (r *Responder) InternalError() error {
	return r.send(http.StatusInternalServerError, defaultInternalErrorResponseMessage)
}

func (r *Responder) send(status int, responseMessage string) error {
	return r.ctx.JSON(status, serverResponse{Response: textResponse{Text: responseMessage}})
}
