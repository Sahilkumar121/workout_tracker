package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// database creation using pgxpool
func Connect(databaseUrl string) (*pgxpool.Pool, error) {

	// pool creation
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	// config
	config.MaxConns = 25

	config.MaxConnIdleTime = time.Minute * 30
	config.MaxConnLifetime = time.Hour * 1

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// check database is alive
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
