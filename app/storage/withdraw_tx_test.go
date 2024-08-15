package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type WithdrawTransactionTestSuite struct {
	suite.Suite
	ctx   context.Context
	pool  *fakePool
	tx    *WithdrawTx
	spyTx *spyTx
}

func Test_WithdrawTransaction(t *testing.T) {
	suite.Run(t, &WithdrawTransactionTestSuite{})
}

func (s *WithdrawTransactionTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.tx = NewWithdrawTx(s.ctx, nil, "")
	s.pool = &fakePool{}
	s.spyTx = &spyTx{}
	s.tx.pool = s.pool
	s.tx.tx = s.spyTx
}

func (s *WithdrawTransactionTestSuite) Test_rowsClossed() {
	s.spyTx.rows = &spyRows{noNext: true}
	_, _ = s.tx.AvailablePoints()
	s.Require().True(s.spyTx.rows.IsClosed)
}

func (s *WithdrawTransactionTestSuite) Test_QueryError() {
	s.spyTx.queryErr = fmt.Errorf("db error")
	_, _ = s.tx.AvailablePoints()
}
