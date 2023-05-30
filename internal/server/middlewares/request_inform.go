package middlewares

import (
	"github.com/labstack/echo"
	logger "github.com/rs/zerolog/log"
	"time"
)

func ApplyRequestInform(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.With().Logger()
		start := time.Now()

		log.Info().
			Str("from", c.Request().Header.Get("Referer")).
			Str("content-type", c.Request().Header.Get("Content-Type")).
			Str("access", c.Request().Header.Get("Access")).
			Str("uri", c.Request().URL.Path).
			Str("request_method", c.Request().Method).Msg("request information")

		err := next(c)

		log.Info().
			Str("uri", c.Request().URL.Path).
			Str("method", c.Request().Method).
			Str("duration", time.Since(start).String()).
			Int64("size", c.Response().Size).
			Int("status", c.Response().Status).
			Msg("request was handled")

		return err
	}
}
