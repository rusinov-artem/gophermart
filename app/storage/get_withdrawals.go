package storage

import (
	"fmt"
	"time"

	"github.com/rusinov-artem/gophermart/app/dto"
)

func (r *Storage) GetWithdrawals(login string) ([]dto.Withdrawal, error) {
	sqlStr := `SELECT order_nr, sum, created_at FROM withdraw WHERE login = $1`

	rows, err := r.pool.Query(r.ctx, sqlStr, login)
	if err != nil {
		return nil, fmt.Errorf("unable to get withdrawals: %w", err)
	}

	defer rows.Close()

	type dbWithdraw struct {
		OrderNr   string
		Sum       float32
		CreatedAt time.Time
	}

	withdraws := []dto.Withdrawal{}

	for rows.Next() {
		w := dbWithdraw{}
		err := rows.Scan(&w.OrderNr, &w.Sum, &w.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to get withdrawals: %w", err)
		}

		withdraws = append(withdraws, dto.Withdrawal{
			OrderNr:     w.OrderNr,
			Sum:         w.Sum,
			ProcessedAt: w.CreatedAt,
		})
	}

	return withdraws, nil

}
