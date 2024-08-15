package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
	"github.com/rusinov-artem/gophermart/app/storage"
	"github.com/rusinov-artem/gophermart/test"
)

type WithdrawTransactionTestSuite struct {
	suite.Suite
	ctx           context.Context
	pool          *pgxpool.Pool
	tx            *storage.WithdrawTx
	storage       *storage.RegistrationStorage
	login         string
	initialPoints float32
	tx2           *storage.WithdrawTx
	cancelFn      context.CancelFunc
	orderNr       string
}

func Test_WithdrawTransaction(t *testing.T) {
	suite.Run(t, &WithdrawTransactionTestSuite{})
}

func (s *WithdrawTransactionTestSuite) SetupSuite() {
	var err error
	dsn := test.CreateTestDB("test_withdraw_transaction")
	s.pool, err = pgxpool.New(context.Background(), dsn)
	s.Require().NoError(err)
}

func (s *WithdrawTransactionTestSuite) SetupTest() {
	s.ctx, s.cancelFn = context.WithTimeout(context.Background(), time.Second)

	s.login = fmt.Sprintf("login_%s", s.T().Name())
	s.orderNr = fmt.Sprintf("order_nr_%s", s.T().Name())

	s.initialPoints = float32(30000)
	s.tx = storage.NewWithdrawTx(s.ctx, s.pool, s.login)
	s.tx2 = storage.NewWithdrawTx(s.ctx, s.pool, s.login)
	s.storage = storage.NewRegistrationStorage(s.ctx, s.pool)
	s.Require().NoError(s.storage.SaveUser(s.login, "password"))

	s.Require().NoError(s.storage.AddOrder(s.login, s.orderNr))

	s.Require().NoError(s.storage.UpdateOrdersState([]dto.OrderListItem{
		{
			OrderNr: s.orderNr,
			Status:  order.PROCESSED,
			Accrual: s.initialPoints,
		},
	}))
}

func (s *WithdrawTransactionTestSuite) TearDownTest() {
	s.cancelFn()
}

func (s *WithdrawTransactionTestSuite) Test_CanWithdraw() {
	withdrawals, err := s.storage.GetWithdrawals(s.login)
	s.Len(withdrawals, 0)
	s.NoError(err)

	s.Require().NoError(s.tx.Begin())
	defer s.tx.Rollback()

	s.Require().NoError(s.tx.LockUser())

	available, err := s.tx.AvailablePoints()
	s.Require().NoError(err)
	s.InDelta(s.initialPoints, available, 0.000001)

	sum := float32(100)
	s.Require().NoError(s.tx.Withdraw(s.orderNr, sum))
	s.Require().NoError(s.tx.Commit())

	withdrawals, err = s.storage.GetWithdrawals(s.login)
	s.NoError(err)
	s.Len(withdrawals, 1)
	s.Equal(s.orderNr, withdrawals[0].OrderNr)
	s.InDelta(sum, withdrawals[0].Sum, 0.0001)
	s.NotEmpty(withdrawals[0].ProcessedAt)

	withdrawn, err := s.storage.Withdrawn(s.login)
	s.NoError(err)
	s.InDelta(sum, withdrawn, 0.0001)

	balance, err := s.storage.Balance(s.login)
	s.NoError(err)
	s.InDelta(s.initialPoints-sum, balance, 0.0001)
}

func (s *WithdrawTransactionTestSuite) Test_CanBlock() {
	s.Require().NoError(s.tx.Begin())

	defer s.tx.Rollback()

	s.Require().NoError(s.tx2.Begin())
	defer s.tx2.Rollback()

	s.Require().NoError(s.tx.LockUser())
	s.Require().ErrorContains(s.tx2.LockUser(), "timeout")
}
