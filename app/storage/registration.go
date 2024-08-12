package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rusinov-artem/gophermart/app/crypto"
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

func (r *RegistrationStorage) IsLoginExists(login string) (bool, error) {
	sqlStr := `SELECT login FROM "user" WHERE login = $1`
	rows, err := r.pool.Query(r.ctx, sqlStr, login)
	if err != nil {
		return false, fmt.Errorf("unable to find user: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	var found string
	err = rows.Scan(&found)
	if err != nil {
		return false, fmt.Errorf("unable to scan: %w", err)
	}

	return true, nil
}

func (r *RegistrationStorage) SaveUser(login, password string) error {
	sqlStr := `INSERT INTO "user" (login, password_hash) VALUES ($1, $2)`
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return fmt.Errorf("unable to hash password: %w", err)
	}

	_, err = r.pool.Exec(r.ctx, sqlStr, login, hash)
	if err != nil {
		return fmt.Errorf("unable to find user: %w", err)
	}

	return nil
}

func (r *RegistrationStorage) AddToken(login, token string) error {
	sqlStr := `INSERT INTO "auth_token" (login, token) VALUES ($1, $2)`

	_, err := r.pool.Exec(r.ctx, sqlStr, login, token)
	if err != nil {
		return fmt.Errorf("unable to add auth_token: %w", err)
	}

	return nil
}
