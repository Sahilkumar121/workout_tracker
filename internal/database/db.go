package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PostgresqlDB(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {

	c, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("Parse Config failed: %w\n", err)
	}

	c.MaxConns = 25
	c.MaxConnIdleTime = 30 * time.Minute
	c.MaxConnLifetime = 1 * time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("creation of database failed: %w\n", err)
	}

	checkCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(checkCtx); err != nil {
		return nil, fmt.Errorf("connection to database failed: %w\n", err)
	}

	return pool, nil
}
