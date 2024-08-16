package storage

import "fmt"

func (r *Storage) FindToken(token string) (string, error) {
	sqlStr := "SELECT login FROM auth_token WHERE token = $1"

	rows, err := r.pool.Query(r.ctx, sqlStr, token)
	if err != nil {
		return "", fmt.Errorf("unable to find token: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return "", fmt.Errorf("token not found: %w", err)
	}

	var login string
	if err := rows.Scan(&login); err != nil {
		return "", fmt.Errorf("unable to find token: %w", err)
	}

	return login, nil
}
