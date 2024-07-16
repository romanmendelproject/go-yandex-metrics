package migrations

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitTables, downInitTables)
}

func upInitTables(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	log.Info("Create DB Table")

	query := `
		CREATE TABLE IF NOT EXISTS metrics (
		    id SERIAL,
		    type VARCHAR(64) NOT NULL,
		    name VARCHAR(128) NOT NULL,
		    counter BIGINT,
		    gauge DOUBLE PRECISION
		)
	`

	// Creating metrics table
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downInitTables(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	log.Info("Remove DB Table")

	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS metrics"); err != nil {
		return err
	}

	return nil
}