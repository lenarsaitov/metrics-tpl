package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lenarsaitov/metrics-tpl/internal/server/models"
	"github.com/rs/zerolog/log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultTimout = time.Second

type PollStorage struct {
	db *sql.DB
}

func NewPollStorage(ctx context.Context, dataSourceName string) (*PollStorage, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %s", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, createCounterMetricsTable)
	if err != nil {
		return nil, fmt.Errorf("failed to create counter metrics table: %s", err)
	}

	_, err = db.ExecContext(ctx, createGaugeMetricsTable)
	if err != nil {
		return nil, fmt.Errorf("failed to create gauge metrics table: %s", err)
	}

	return &PollStorage{db}, nil
}

func (m *PollStorage) GetAll(ctx context.Context) (models.Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimout)
	defer cancel()

	metrics := models.Metrics{
		GaugeMetrics:   make([]models.GaugeMetric, 0),
		CounterMetrics: make([]models.CounterMetric, 0),
	}

	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return metrics, err
	}
	defer func() {
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Error().Err(errRollback).Msg("failed to rollback transaction")
			}
		} else {
			errCommit := tx.Commit()
			if errCommit != nil {
				log.Error().Err(errCommit).Msg("failed to commit transaction")
			}
		}
	}()

	rows, err := tx.QueryContext(ctx, "SELECT name, value FROM gauge_metrics")
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		metric := models.GaugeMetric{}
		err = rows.Scan(&metric.Name, &metric.Value)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan fields from gauge metrics table")

			return metrics, err
		}

		metrics.GaugeMetrics = append(metrics.GaugeMetrics, metric)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("error when parse rows")

		return metrics, err
	}

	rows, err = tx.QueryContext(ctx, "SELECT name, value FROM counter_metrics")
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		metric := models.CounterMetric{}
		err = rows.Scan(&metric.Name, &metric.Value)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan fields from counter metrics table")

			return metrics, err
		}

		metrics.CounterMetrics = append(metrics.CounterMetrics, metric)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("error when parse rows")

		return metrics, err
	}

	return metrics, nil
}

func (m *PollStorage) GetGaugeMetric(ctx context.Context, name string) (*float64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimout)
	defer cancel()

	row := m.db.QueryRowContext(ctx, "SELECT value FROM gauge_metrics WHERE name = $1", name)

	var value float64
	err := row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to get gauge metric")

		return nil, err
	}

	return &value, nil
}

func (m *PollStorage) GetCounterMetric(ctx context.Context, name string) (*int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimout)
	defer cancel()

	row := m.db.QueryRowContext(ctx, "SELECT value FROM counter_metrics WHERE name = $1", name)

	var value int64
	err := row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to get gauge metric")

		return nil, err
	}

	return &value, nil
}

func (m *PollStorage) ReplaceGauge(ctx context.Context, name string, value float64) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimout)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")

		return err
	}
	defer func() {
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Error().Err(errRollback).Msg("failed to rollback transaction")
			}
		} else {
			errCommit := tx.Commit()
			if errCommit != nil {
				log.Error().Err(errCommit).Msg("failed to commit transaction")
			}
		}
	}()

	var previousValue float64
	err = tx.QueryRowContext(ctx, "SELECT value FROM gauge_metrics WHERE name = $1 FOR UPDATE", name).Scan(&previousValue)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error().Err(err).Msg("failed to get gauge metric (but not errNoRows)")

		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		_, err = tx.ExecContext(ctx, "INSERT INTO gauge_metrics (name, value) VALUES ($1, $2)", name, value)
		if err != nil {
			log.Error().Err(err).Msg("failed to create gauge metric")

			return err
		}

		return nil
	}

	_, err = tx.ExecContext(ctx, "UPDATE gauge_metrics SET value = $1, updated_at = $2 WHERE name = $3", value, time.Now(), name)
	if err != nil {
		log.Error().Err(err).Msg("failed to exec request update gauge metric")

		return err
	}

	return nil
}

func (m *PollStorage) AddCounter(ctx context.Context, name string, delta int64) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimout)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")

		return 0, err
	}
	defer func() {
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Error().Err(errRollback).Msg("failed to rollback transaction")
			}
		} else {
			errCommit := tx.Commit()
			if errCommit != nil {
				log.Error().Err(errCommit).Msg("failed to commit transaction")
			}
		}
	}()

	var previousValue int64
	err = tx.QueryRowContext(ctx, "SELECT value FROM counter_metrics WHERE name = $1 FOR UPDATE", name).Scan(&previousValue)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error().Err(err).Msg("failed to get counter metric (but not errNoRows)")

		return 0, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		_, err = tx.ExecContext(ctx, "INSERT INTO counter_metrics (name, delta, value) VALUES ($1, $2, $3)", name, delta, delta)
		if err != nil {
			log.Error().Err(err).Msg("failed to create counter metric")

			return 0, err
		}

		return delta, nil
	}

	newValue := previousValue + delta
	_, err = tx.ExecContext(ctx, "UPDATE counter_metrics SET value = $1, delta = $2, updated_at = $3 WHERE name = $4", newValue, delta, time.Now(), name)
	if err != nil {
		log.Error().Err(err).Msg("failed to exec request update counter metric")

		return 0, err
	}

	return newValue, nil
}
