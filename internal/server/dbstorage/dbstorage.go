// Модуль для взаимодействия с БД
package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/metrics"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"
	log "github.com/sirupsen/logrus"
)

// PostgresStorage определяет объект для работы с БД
type PostgresStorage struct {
	db *pgxpool.Pool
}

var (
	pgInstance *PostgresStorage
	pgOnce     sync.Once
)

// NewPostgresStorage создает объект для работы с БД
func NewPostgresStorage(ctx context.Context, connString string) *PostgresStorage {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Fatal("unable to create connection pool: %w", err)
		}

		pgInstance = &PostgresStorage{db}
	})

	return pgInstance
}

// Ping проверяет доступность БД
func (pg *PostgresStorage) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

// Close закрывает соединение с БД
func (pg *PostgresStorage) Close() {
	pg.db.Close()
}

// SetGauge записывает данные формата Gauge в БД
func (pg *PostgresStorage) SetGauge(ctx context.Context, name string, value float64) error {
	var oldVal float64

	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Check if metric exists
	if err := tx.QueryRow(ctx, `SELECT gauge FROM metrics WHERE name=$1 AND type = 'gauge'`, name).Scan(&oldVal); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Insert new metric if not exists
			log.Error(name)
			if _, err := tx.Exec(ctx, `INSERT INTO metrics (name, type, gauge) VALUES ($1, 'gauge', $2)`, name, value); err != nil {
				log.Error(err)
				return err
			}
			return nil
		}

		log.Error(err)
		return err
	}

	// Update metric if exists
	if _, err := tx.Exec(ctx, `UPDATE metrics SET gauge = $1 WHERE type = 'gauge' AND name = $2`, value, name); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// SetCounter записывает данные формата Counter в БД
func (pg *PostgresStorage) SetCounter(ctx context.Context, name string, value int64) error {
	var oldVal int64

	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	if err := tx.QueryRow(ctx, "SELECT counter FROM metrics WHERE name = $1 AND type = 'counter'", name).Scan(&oldVal); err != nil {
		// Insert new metric if not exists
		if errors.Is(err, pgx.ErrNoRows) {
			if _, err := tx.Exec(ctx, `INSERT INTO metrics (name, type, counter) VALUES ($1, 'counter', $2)`, name, value); err != nil {
				log.Error(err)
				return err
			}

			return nil
		}
		return err
	}

	// Update metric if exists
	value += oldVal
	if _, err := tx.Exec(ctx, `UPDATE metrics SET counter = $1 WHERE type = 'counter' AND name = $2`, value, name); err != nil {
		log.Error(err)
		return err
	}

	return nil

}

// GetCounter читает данные формата Counter из БД
func (pg *PostgresStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	var counter sql.NullInt64

	if err := pg.db.QueryRow(ctx, "SELECT counter FROM metrics WHERE name = $1 AND type = 'counter'", name).Scan(&counter); err != nil {
		return 0, err
	}

	if !counter.Valid {
		return 0, fmt.Errorf("unexpected type of metric")
	}

	return counter.Int64, nil
}

// GetGauge читает данные формата Gauge из БД
func (pg *PostgresStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	var gauge sql.NullFloat64

	if err := pg.db.QueryRow(ctx, "SELECT gauge FROM metrics WHERE name = $1 AND type = 'gauge'", name).Scan(&gauge); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error(err)
			return 0, nil
		}
		return 0, err
	}

	if !gauge.Valid {
		return 0, fmt.Errorf("unexpected type of metric")
	}

	return gauge.Float64, nil
}

// GetGauge все данные из БД
func (pg *PostgresStorage) GetAll(ctx context.Context) ([]storage.Value, error) {
	var values []storage.Value

	rows, err := pg.db.Query(ctx, `SELECT type, name, gauge, counter FROM metrics`)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mType string
		var name string
		var gauge sql.NullFloat64
		var counter sql.NullInt64

		if err := rows.Scan(&mType, &name, &gauge, &counter); err != nil {
			log.Error(err)
			return nil, err
		}

		if gauge.Valid {
			values = append(values, storage.Value{
				Name:  name,
				Type:  "gauge",
				Value: strconv.FormatFloat(gauge.Float64, 'f', 1, 64),
			})
		} else if counter.Valid {
			values = append(values, storage.Value{
				Name:  name,
				Type:  "counter",
				Value: counter.Int64,
			})
		} else {
			log.Error(err)
			return nil, err
		}
	}

	return values, nil
}

// SetBatch записывает данные в БД с использование одного запроса
func (pg *PostgresStorage) SetBatch(ctx context.Context, metrics []metrics.Metric) error {
	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	batch := &pgx.Batch{}

	query := `INSERT INTO metrics (type, name, counter, gauge) VALUES ($1, $2, $3, $4)`

	for _, metric := range metrics {
		if metric.MType == "counter" {
			var oldCounter int64
			if err := tx.QueryRow(ctx, `SELECT counter FROM metrics WHERE name=$1 AND type = 'counter'`, metric.ID).Scan(&oldCounter); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					if _, err := tx.Exec(ctx, query, metric.MType, metric.ID, metric.Delta, metric.Value); err != nil {
						log.Error(err)
						return err
					}
					continue
				}

				log.Error(err)
				return err
			}

			*metric.Delta += oldCounter

			batch.Queue(`UPDATE metrics SET counter = $1 WHERE type = 'counter' AND name = $2`, metric.Delta, metric.ID)

		} else if metric.MType == "gauge" {
			var oldGauge float64
			if err := tx.QueryRow(ctx, `SELECT 1 FROM metrics WHERE name=$1 AND type = 'gauge'`, metric.ID).Scan(&oldGauge); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					if _, err := tx.Exec(ctx, query, metric.MType, metric.ID, metric.Delta, metric.Value); err != nil {
						log.Error(err)
						return err
					}
					continue
				}

				log.Error(err)
				return err
			}

			batch.Queue(`UPDATE metrics SET gauge = $1 WHERE type = 'gauge' AND name = $2`, metric.Value, metric.ID)
		}
	}

	br := tx.SendBatch(ctx, batch)
	br.Close()

	return nil
}
