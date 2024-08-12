package storage

import "fmt"

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
