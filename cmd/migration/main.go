package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxstd "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/ninedraft/core-api/migrations"
	"log/slog"
	"os"
	"time"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	godotenv.Load()

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		panic("POSTGRES_DSN is required")
	}

	ctx := context.Background()

	slog.DebugContext(ctx, "connecting to postgres")
	sqlPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}
	defer sqlPool.Close()

	go func() {
		for range time.Tick(4 * time.Second) {
			if err := sqlPool.Ping(ctx); err != nil {
				slog.ErrorContext(ctx, "failed to ping postgres", "error", err)
			}
		}
	}()

	db := pgxstd.OpenDBFromPool(sqlPool)

	if err = migrations.Up(ctx, db); err != nil {
		panic("migrations failed" + err.Error())
	}
}
