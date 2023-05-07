package usecase

import (
	"github.com/rs/zerolog"
)

type (
	MetricsAgent interface {
		PollAndReport(log *zerolog.Logger)
	}
)
