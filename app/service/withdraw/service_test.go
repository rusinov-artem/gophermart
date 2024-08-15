package withdraw

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type WithdrawServiceTestSuite struct {
	suite.Suite
	tx      *tx
	service *Service
	logger  *zap.Logger
	logs    *bytes.Buffer
}

func Test_WithdrawService(t *testing.T) {
	suite.Run(t, &WithdrawServiceTestSuite{})
}

func (s *WithdrawServiceTestSuite) SetupTest() {
	s.tx = &tx{}

	txFactory := func(login string) Transaction {
		s.tx.login = login
		return s.tx
	}

	s.logger, s.logs = logger.SpyLogger()

	s.service = NewWithdrawService(txFactory, s.logger)
}

func (s *WithdrawServiceTestSuite) Test_UnableToBeginTx() {
	params := dto.WithdrawParams{Login: "login"}
	s.tx.beginErr = fmt.Errorf("database error")
	err := s.service.Withdraw(params)
	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(params.Login, s.tx.login)
	s.Contains(s.logs.String(), s.tx.beginErr.Error())
}

func (s *WithdrawServiceTestSuite) Test_UnableToLockUser() {
	params := dto.WithdrawParams{Login: "login"}
	s.tx.lockUserErr = fmt.Errorf("network error")
	err := s.service.Withdraw(params)
	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(params.Login, s.tx.login)
	s.Contains(s.logs.String(), s.tx.lockUserErr.Error())
	s.True(s.tx.rollbackExecuted)
}

func (s *WithdrawServiceTestSuite) Test_UnableToGetAvailablePoints() {
	params := dto.WithdrawParams{Login: "login"}
	s.tx.availablePointsErr = fmt.Errorf("connection refused")
	err := s.service.Withdraw(params)
	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(params.Login, s.tx.login)
	s.Contains(s.logs.String(), s.tx.availablePointsErr.Error())
	s.True(s.tx.rollbackExecuted)
}

func (s *WithdrawServiceTestSuite) Test_NotEnoughPoints() {
	params := dto.WithdrawParams{Login: "login", Sum: 999}
	s.tx.availablePoints = 100

	err := s.service.Withdraw(params)
	s.NotNil(err)
	s.Equal(appError.NotEnoughPoints, err.Code)
	s.Equal(params.Login, s.tx.login)
	s.True(s.tx.rollbackExecuted)
}

func (s *WithdrawServiceTestSuite) Test_UnableToWithdraw() {
	params := dto.WithdrawParams{Login: "login", Sum: 999, OrderNr: "orderNr"}
	s.tx.availablePoints = 10000

	s.tx.withdrawErr = fmt.Errorf("db error")

	err := s.service.Withdraw(params)
	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(params.Login, s.tx.login)
	s.Equal(params.OrderNr, s.tx.orderNr)
	s.InDelta(params.Sum, s.tx.amount, 0.001)
	s.True(s.tx.rollbackExecuted)
	s.Contains(s.logs.String(), s.tx.withdrawErr.Error())
}

func (s *WithdrawServiceTestSuite) Test_UnableToCommit() {
	params := dto.WithdrawParams{Login: "login", Sum: 999, OrderNr: "orderNr"}
	s.tx.availablePoints = 10000

	s.tx.commitErr = fmt.Errorf("commit error")

	err := s.service.Withdraw(params)
	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(params.Login, s.tx.login)
	s.Equal(params.OrderNr, s.tx.orderNr)
	s.InDelta(params.Sum, s.tx.amount, 0.001)
	s.True(s.tx.rollbackExecuted)
	s.Contains(s.logs.String(), s.tx.commitErr.Error())
}

func (s *WithdrawServiceTestSuite) Test_Success() {
	params := dto.WithdrawParams{Login: "login", Sum: 999, OrderNr: "orderNr"}
	s.tx.availablePoints = 10000

	err := s.service.Withdraw(params)
	s.Nil(err)
	s.Equal(params.Login, s.tx.login)
	s.Equal(params.OrderNr, s.tx.orderNr)
	s.InDelta(params.Sum, s.tx.amount, 0.001)
	s.True(s.tx.rollbackExecuted)
}

type tx struct {
	login    string
	beginErr error

	rollbackExecuted bool
	lockUserErr      error

	availablePointsErr error
	availablePoints    float32

	withdrawErr error
	orderNr     string
	amount      float32

	commited  bool
	commitErr error
}

func (t *tx) Begin() error {
	return t.beginErr
}

func (t *tx) Rollback() {
	t.rollbackExecuted = true
}

func (t *tx) LockUser() error {
	return t.lockUserErr
}

func (t *tx) AvailablePoints() (float32, error) {
	return t.availablePoints, t.availablePointsErr
}

func (t *tx) Withdraw(orderNr string, amount float32) error {
	t.orderNr = orderNr
	t.amount = amount
	return t.withdrawErr
}

func (t *tx) Commit() error {
	t.commited = true
	return t.commitErr
}
