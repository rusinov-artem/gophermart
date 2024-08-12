package storage

import "fmt"

func (r *RegistrationStorage) AddOrder(login, orderNr string) error {
	sqlStr := `INSERT INTO "order" ( login, order_nr) VALUES ( $1, $2 )`

	_, err := r.pool.Exec(r.ctx, sqlStr, login, orderNr)
	if err != nil {
		return fmt.Errorf("unable to add order: %w", err)
	}

	return nil
}
