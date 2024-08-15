package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type querier interface {
	Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error)
}

func getBalance(ctx context.Context, q querier, login string) (float32, error) {
	sqlStr := `
	SELECT 
    COALESCE((SELECT sum(accrual) FROM "order" WHERE login = $1 AND status = 'PROCESSED'), 0) -
    COALESCE((SELECT sum(withdraw.sum)   FROM "withdraw" WHERE login = $2), 0) as balance`

	rows, err := q.Query(ctx, sqlStr, login, login)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	if !rows.Next() {
		return 0, err
	}

	var available float32
	err = rows.Scan(&available)

	return available, err
}
