package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pgxstd "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/ninedraft/core-api/internal/api"
	"github.com/ninedraft/core-api/internal/service"
	"github.com/ninedraft/core-api/internal/storage"
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

	st := storage.New(db)
	srv := service.New(st)

	strict := api.NewStrictHandlerWithOptions(srv, nil,
		api.StrictHTTPServerOptions{
			RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				slog.ErrorContext(r.Context(), "request error", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
			},
			ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				slog.ErrorContext(r.Context(), "request error", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			},
		})
	handler := api.Handler(strict)

	http.Handle("/", logMw(handler))

	slog.DebugContext(ctx, "listening on :8081")
	if err := http.ListenAndServe("localhost:8081", nil); err != nil {
		panic(err)
	}
}

func logMw(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ri := &responseInterceptor{ResponseWriter: w}
		defer func() {
			slog.DebugContext(r.Context(), r.Method,
				"url", r.URL.String(),
				"status", ri.status)

		}()

		next.ServeHTTP(ri, r)
	}
}

type responseInterceptor struct {
	http.ResponseWriter

	status int
}

func (ri *responseInterceptor) WriteHeader(status int) {
	ri.status = status
	ri.ResponseWriter.WriteHeader(status)
}

func (ri *responseInterceptor) Write(b []byte) (int, error) {
	if ri.status == 0 {
		ri.status = http.StatusOK
	}

	return ri.ResponseWriter.Write(b)
}
