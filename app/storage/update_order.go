package storage

import (
	"github.com/jackc/pgx/v5"

	"github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
)

func (r *Storage) UpdateOrdersState(orders []dto.OrderListItem) error {
	ProcessedSQL := `UPDATE "order" SET status = $1, accrual = $2 WHERE order_nr = $3`
	progressSQL := `UPDATE "order" SET status = $1 WHERE order_nr = $2`

	batch := pgx.Batch{}

	for i := range orders {
		if orders[i].Status == order.PROCESSED {
			batch.Queue(ProcessedSQL, orders[i].Status, orders[i].Accrual, orders[i].OrderNr)
		} else {
			batch.Queue(progressSQL, orders[i].Status, orders[i].OrderNr)
		}
	}

	batchResult := r.pool.SendBatch(r.ctx, &batch)
	return batchResult.Close()
}
