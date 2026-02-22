package postgres

import (
	"context"
	"errors"
	"fivestars/internal/infra/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a PostgreSQL connection pool from dsn (e.g. DATABASE_URL).
func NewPool(ctx context.Context, pc config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := pc.DSN()
	if dsn == "" {
		return nil, errors.New("Database DSN is empty")
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.New("parse config error: " + err.Error())
	}

	poolConfig.MaxConns = pc.MaxConns
	poolConfig.MinConns = pc.MinConns
	poolConfig.HealthCheckPeriod = 1 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Timeout de conexão
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
