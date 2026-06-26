package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func Connect(ctx context.Context, url string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to database", "error", err)
		os.Exit(1)
	}

	return db, nil
}

func RunMigrations(ctx context.Context, db *sqlx.DB) error {
	slog.DebugContext(ctx, "Running database migrations...")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create database driver", "error", err)
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver,
	)

	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	slog.Info("Database migrations completed successfully")
	return nil
}
