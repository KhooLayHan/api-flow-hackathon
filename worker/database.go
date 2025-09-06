package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newDBConnectionPool(logger *slog.Logger) (*pgxpool.Pool, error) {
	connPool, connPoolErr := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if connPoolErr != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", connPoolErr)
	}

	if err := connPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection pool created successfully.")
	return connPool, nil
}
