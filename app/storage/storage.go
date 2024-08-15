package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
}

type RegistrationStorage struct {
	pool pool
	ctx  context.Context
}

func NewRegistrationStorage(ctx context.Context, pool *pgxpool.Pool) *RegistrationStorage {
	return &RegistrationStorage{
		pool: pool,
		ctx:  ctx,
	}
}
