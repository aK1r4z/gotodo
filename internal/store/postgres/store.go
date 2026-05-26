package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dbURL string) (*db, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	return &db{
		pool: pool,
	}, nil
}

func (d *db) Close() {
	d.pool.Close()
}
