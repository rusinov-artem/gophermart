package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rusinov-artem/gophermart/app/dto"
)

func (r *RegistrationStorage) ListOrders(login string) ([]dto.OrderListItem, error) {
	sqlStr := `SELECT order_nr, status, accrual, upload_at
				FROM "order" 
				WHERE login = $1
				ORDER BY upload_at ASC
`

	rows, err := r.pool.Query(r.ctx, sqlStr, login)
	if err != nil {
		return nil, fmt.Errorf("unable to list orders: %w", err)
	}
	defer rows.Close()

	type dbOrder struct {
		OrderNr  string
		Status   string
		Accrual  sql.NullInt64
		UploadAt time.Time
	}

	orders := []dto.OrderListItem{}

	for rows.Next() {
		var o dbOrder
		err := rows.Scan(&o.OrderNr, &o.Status, &o.Accrual, &o.UploadAt)
		if err != nil {
			return nil, fmt.Errorf("unable to list orders: %w", err)
		}
		orders = append(orders, dto.OrderListItem{
			OrderNr:  o.OrderNr,
			Status:   o.Status,
			Accrual:  o.Accrual.Int64,
			UploadAt: o.UploadAt,
		})
	}

	return orders, nil
}
