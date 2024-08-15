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
	s.Require().NoError(s.tx.Begin())
	defer s.tx.Rollback()

	s.Require().NoError(s.tx.LockUser())

	available, err := s.tx.AvailablePoints()
	s.Require().NoError(err)
	s.InDelta(s.initialPoints, available, 0.000001)

	s.Require().NoError(s.tx.Withdraw(s.orderNr, 100))
	s.Require().NoError(s.tx.Commit())
}

func (s *WithdrawTransactionTestSuite) Test_CanBlock() {
	s.Require().NoError(s.tx.Begin())

	defer s.tx.Rollback()

	s.Require().NoError(s.tx2.Begin())
	defer s.tx2.Rollback()

	s.Require().NoError(s.tx.LockUser())
	s.Require().ErrorContains(s.tx2.LockUser(), "timeout")
}
