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

type DBStorage struct {
	db *pgxpool.Pool
}

var (
	pgInstance *DBStorage
	pgOnce     sync.Once
)

func NewDBStorage(ctx context.Context, connString string) *DBStorage {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			panic(fmt.Errorf("unable to create connection pool: %w", err))
		}

		pgInstance = &DBStorage{db}
	})

	return pgInstance
}

func (pg *DBStorage) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *DBStorage) Close() {
	pg.db.Close()
}

func (pg *DBStorage) SetGauge(ctx context.Context, name string, value float64) error {
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

func (pg *DBStorage) SetCounter(ctx context.Context, name string, value int64) error {
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

func (pg *DBStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	var counter sql.NullInt64

	if err := pg.db.QueryRow(ctx, "SELECT counter FROM metrics WHERE name = $1 AND type = 'counter'", name).Scan(&counter); err != nil {
		return 0, err
	}

	if !counter.Valid {
		return 0, fmt.Errorf("unexpected type of metric")
	}

	return counter.Int64, nil
}

func (pg *DBStorage) GetGauge(ctx context.Context, name string) (float64, error) {
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

func (pg *DBStorage) GetAll(ctx context.Context) ([]storage.Value, error) {
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

func (pg *DBStorage) SetBatch(ctx context.Context, metrics []metrics.Metric) error {
	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

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

			if _, err := tx.Exec(ctx, `UPDATE metrics SET counter = $1 WHERE type = 'counter' AND name = $2`, metric.Delta, metric.ID); err != nil {
				log.Error(err)
				return err
			}
		} else if metric.MType == "gauge" {
			ra, err := tx.Exec(ctx, `UPDATE metrics SET gauge = $1 WHERE type = 'gauge' AND name = $2`, metric.Value, metric.ID)
			if err != nil {
				log.Error(err)
				return err
			}

			if ra.RowsAffected() == 0 {
				if _, err := tx.Exec(ctx, query, metric.MType, metric.ID, metric.Delta, metric.Value); err != nil {
					log.Error(err)
					return err
				}
			}
		}
	}

	return nil
}
