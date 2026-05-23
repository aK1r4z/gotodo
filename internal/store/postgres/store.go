package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type db struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dbURL string) (*db, error) {
	pool, err := pgxpool.Connect(ctx, dbURL)
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
