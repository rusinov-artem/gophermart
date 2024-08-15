package get

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

type GetBalanceActionTestSuite struct {
	suite.Suite
	action       *Action
	orderService *orderService
	storage      *storage
	logs         *bytes.Buffer
	logger       *zap.Logger
}

func Test_GetBalanceAction(t *testing.T) {
	suite.Run(t, &GetBalanceActionTestSuite{})
}

func (s *GetBalanceActionTestSuite) SetupTest() {
	s.orderService = &orderService{}
	s.storage = &storage{}
	s.logger, s.logs = logger.SpyLogger()
	s.action = New(s.orderService, s.storage, s.logger)
}

func (s *GetBalanceActionTestSuite) Test_UnableToGetUserOrders() {
	login := "login"
	s.orderService.err = &appError.InternalError{
		InnerError: fmt.Errorf("database error"),
		Msg:        "service temporary unavailable",
		Code:       appError.ServiceUnavailable,
	}

	_, appErr := s.action.GetBalance(login)
	s.NotNil(appErr)
	s.Equal(login, s.orderService.ordersForLogin)
	s.Equal(s.orderService.err, appErr)
}

func (s *GetBalanceActionTestSuite) Test_unableToGetBalance() {
	login := "login"
	s.storage.balanceErr = fmt.Errorf("database error")
	_, appErr := s.action.GetBalance(login)
	s.NotNil(appErr)
	s.Equal(login, s.storage.loginForBalance)
	s.Contains(s.logs.String(), s.storage.balanceErr.Error())
}

func (s *GetBalanceActionTestSuite) Test_unableToGetWithdrawn() {
	login := "login"
	s.storage.balance = 42.42
	s.storage.withdrawnErr = fmt.Errorf("database error")
	_, appErr := s.action.GetBalance(login)
	s.NotNil(appErr)
	s.Equal(login, s.storage.loginForBalance)
	s.Equal(login, s.storage.loginForWithdrawn)
	s.Contains(s.logs.String(), s.storage.withdrawnErr.Error())
}

func (s *GetBalanceActionTestSuite) Test_Success() {
	login := "login"
	s.storage.balance = 42.42
	s.storage.withdrawn = 55.55
	res, appErr := s.action.GetBalance(login)
	s.Nil(appErr)
	s.Equal(login, s.storage.loginForBalance)
	s.Equal(login, s.storage.loginForWithdrawn)
	s.InDelta(s.storage.balance, res.Current, 0.001)
	s.InDelta(s.storage.withdrawn, res.Withdrawn, 0.001)
}

type orderService struct {
	ordersForLogin string
	orders         []dto.OrderListItem
	err            *appError.InternalError
}

func (t *orderService) GetUserOrders(login string) ([]dto.OrderListItem, *appError.InternalError) {
	t.ordersForLogin = login
	return t.orders, t.err
}

type storage struct {
	loginForBalance string
	balance         float32
	balanceErr      error

	loginForWithdrawn string
	withdrawn         float32
	withdrawnErr      error
}

func (s *storage) Balance(login string) (float32, error) {
	s.loginForBalance = login
	return s.balance, s.balanceErr
}

func (s *storage) Withdrawn(login string) (float32, error) {
	s.loginForWithdrawn = login
	return s.withdrawn, s.withdrawnErr
}
