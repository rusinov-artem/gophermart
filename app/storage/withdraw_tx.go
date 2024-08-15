package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WithdrawTx struct {
	ctx   context.Context
	pool  pool
	login string
	tx    pgx.Tx
}

func NewWithdrawTx(ctx context.Context, p *pgxpool.Pool, login string) *WithdrawTx {
	return &WithdrawTx{
		ctx:   ctx,
		pool:  p,
		login: login,
	}
}

func (t *WithdrawTx) Begin() error {
	var err error
	t.tx, err = t.pool.Begin(t.ctx)
	return err
}

func (t *WithdrawTx) Rollback() {
	_ = t.tx.Rollback(t.ctx)
}

func (t *WithdrawTx) LockUser() error {
	sqlStr := `SELECT login FROM "user" WHERE login = $1 FOR UPDATE`

	_, err := t.tx.Exec(t.ctx, sqlStr, t.login)
	return err
}

func (t *WithdrawTx) AvailablePoints() (float32, error) {
	return getBalance(t.ctx, t.tx, t.login)
}

func (t *WithdrawTx) Withdraw(orderNr string, amount float32) error {
	sqlStr := `INSERT INTO withdraw (login, order_nr, sum) VALUES($1, $2, $3)`
	_, err := t.tx.Exec(t.ctx, sqlStr, t.login, orderNr, amount)
	return err
}

func (t *WithdrawTx) Commit() error {
	return t.tx.Commit(t.ctx)
}
