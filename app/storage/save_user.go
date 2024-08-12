package storage

import (
	"fmt"

	"github.com/rusinov-artem/gophermart/app/crypto"
)

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
