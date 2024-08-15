package storage

import "fmt"

func (r *RegistrationStorage) Withdrawn(login string) (float32, error) {
	sqlStr := `SELECT COALESCE(sum(withdraw.sum),0) as withdrawn FROM withdraw WHERE login = $1;`

	rows, err := r.pool.Query(r.ctx, sqlStr, login)
	if err != nil {
		return 0, fmt.Errorf("unable to get withdrawn: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("unable to get withdrawn: no rows found")
	}

	var withdrawn float32
	err = rows.Scan(&withdrawn)

	return withdrawn, err
}
