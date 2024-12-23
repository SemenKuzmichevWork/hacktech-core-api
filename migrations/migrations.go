package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var fsys embed.FS

func Up(ctx context.Context, db *sql.DB) error {
	migrator, err := newMigrator(db)
	if err != nil {
		return fmt.Errorf("failed to create migrations provider: %w", err)
	}

	results, err := migrator.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	for _, result := range results {
		slog.InfoContext(ctx, "migration applied",
			"migration", result)
	}

	return nil
}

func Reset(ctx context.Context, db *sql.DB) error {
	migrator, err := newMigrator(db)
	if err != nil {
		return fmt.Errorf("failed to create migrations provider: %w", err)
	}

	for {
		result, err := migrator.Down(ctx)
		if errors.Is(err, goose.ErrNoNextVersion) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		slog.InfoContext(ctx, "migration reverted",
			"migration", result)
	}

	return nil
}

func newMigrator(db *sql.DB) (*goose.Provider, error) {
	return goose.NewProvider(goose.DialectPostgres, db, fsys)
}
