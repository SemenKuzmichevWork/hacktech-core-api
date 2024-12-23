package testx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxstd "github.com/jackc/pgx/v5/stdlib"
	"github.com/ninedraft/core-api/migrations"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresUser         = "postgres"
	postgresPassword     = "test-password"
	postgresDatabase     = "test-database"
	postgresInternalPort = 5432
)

func AssertFast(t testing.TB) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
}

func postgresDSN(port nat.Port) string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", postgresUser, postgresPassword, port.Int(), postgresDatabase)
}

var postgres = testcontainers.ContainerRequest{
	Name:  "postgres-test-core-api-" + uuid.NewString(),
	Image: "postgres:17.1-alpine3.20",
	ExposedPorts: []string{
		fmt.Sprintf("%d/tcp", postgresInternalPort),
	},
	WaitingFor: wait.ForAll(
		wait.ForExec([]string{"psql", "-U", postgresUser, "-d", postgresDatabase, "-c", "SELECT 1"}),
		wait.ForExposedPort(),
	),
	Env: map[string]string{
		"POSTGRES_USER":     postgresUser,
		"POSTGRES_PASSWORD": postgresPassword,
		"POSTGRES_DB":       postgresDatabase,
	},
}

func Postgres(t testing.TB) (_ testcontainers.Container, db *sql.DB) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: postgres,
		Started:          true,
		Reuse:            true,
	}

	t.Log("starting postgres container")
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			slog.Error("failed to terminate container", "error", err)
		}
	})

	exposedPort, err := container.MappedPort(ctx, nat.Port(strconv.Itoa(postgresInternalPort)))
	if err != nil {
		t.Fatalf("failed to get exposed port: %v", err)
	}

	dsn := postgresDSN(exposedPort)
	t.Logf("connecting to postgres: %q", dsn)

	sqlPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create pgx connection pool: %v", err)
	}

	db = pgxstd.OpenDBFromPool(sqlPool)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close db connection", "error", err)
		}
	})

	t.Log("applying migrations")
	if err := migrations.Up(ctx, db); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	t.Cleanup(func() {
		if err := migrations.Reset(context.WithoutCancel(ctx), db); err != nil {
			slog.Error("failed to rollback migrations", "error", err)
		}
	})

	return container, db
}
