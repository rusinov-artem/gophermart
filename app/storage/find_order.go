package storage

import (
	"fmt"

	order2 "github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
)

func (r *Storage) FindOrder(orderNr string) (dto.Order, error) {
	var order dto.Order
	sqlStr := `SELECT order_nr, login FROM "order" WHERE order_nr = $1`

	rows, err := r.pool.Query(r.ctx, sqlStr, orderNr)
	if err != nil {
		return order, fmt.Errorf("unable to find order: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return order, &order2.NotFoundErr{OrderNr: orderNr}
	}

	err = rows.Scan(&order.OrderNr, &order.Login)
	if err != nil {
		return order, fmt.Errorf("unable to find order: %w", err)
	}

	return order, nil
}
