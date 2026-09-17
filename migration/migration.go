package migration

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/Ow1Dev/gitria-git/pkgs/zerogoose"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

//go:embed *.sql
var migrationsFS embed.FS

const dialect = "sqlite3"

func RunMigration(db *sql.DB, logger zerolog.Logger) error {
	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(zerogoose.New(logger))

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info().Msg("Migrations completed successfully")
	return nil
}

func MigrationDown(db *sql.DB, logger zerolog.Logger) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(db, "."); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	logger.Info().Msg("Migrations successfully rollbacked")
	return nil
}
