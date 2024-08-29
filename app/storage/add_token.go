package storage

import "fmt"

func (r *Storage) AddToken(login, token string) error {
	sqlStr := `INSERT INTO "auth_token" (login, token) VALUES ($1, $2)`

	_, err := r.pool.Exec(r.ctx, sqlStr, login, token)
	if err != nil {
		return fmt.Errorf("unable to add auth_token: %w", err)
	}

	return nil
}
