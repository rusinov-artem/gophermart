package storage

import (
	"fmt"

	"github.com/rusinov-artem/gophermart/app/crypto"
)

func (r *Storage) SaveUser(login, password string) error {
	sqlStr := `INSERT INTO "user" (login, password_hash) VALUES ($1, $2)`
	_, err := r.pool.Exec(r.ctx, sqlStr, login, crypto.HashPassword(password))
	if err != nil {
		return fmt.Errorf("unable to save user: %w", err)
	}

	return nil
}
