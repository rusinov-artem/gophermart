package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RegistrationStorage struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewRegistrationStorage(ctx context.Context, pool *pgxpool.Pool) *RegistrationStorage {
	return &RegistrationStorage{
		pool: pool,
		ctx:  ctx,
	}
}
