package process

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type WithdrawActionTestSuite struct {
	suite.Suite
	action          *Action
	orderService    *orderService
	withdrawService *withdrawService
}

func Test_WithdrawAction(t *testing.T) {
	suite.Run(t, &WithdrawActionTestSuite{})
}

func (s *WithdrawActionTestSuite) SetupTest() {
	s.orderService = &orderService{}
	s.withdrawService = &withdrawService{}
	s.action = New(s.orderService, s.withdrawService)
}

func (s *WithdrawActionTestSuite) Test_UnableGetBalance() {
	login := "login"
	params := dto.WithdrawParams{
		Login: login,
	}

	s.orderService.err = &appError.InternalError{
		Msg:  "Service temporary unavailable",
		Code: appError.ServiceUnavailable,
	}

	apperr := s.action.Withdraw(params)
	s.NotNil(apperr)
	s.Equal(login, s.orderService.login)
	s.Equal(appError.ServiceUnavailable, apperr.Code)
}

func (s *WithdrawActionTestSuite) Test_UserHasNoOrders() {
	login := "login"
	params := dto.WithdrawParams{
		Login: login,
	}

	s.orderService.err = &appError.InternalError{
		Msg:  "Service temporary unavailable",
		Code: appError.NoOrdersFound,
	}

	apperr := s.action.Withdraw(params)
	s.NotNil(apperr)
	s.Equal(login, s.orderService.login)
	s.Equal(appError.NotEnoughPoints, apperr.Code)
}

func (s *WithdrawActionTestSuite) Test_UnableToWithdraw() {
	login := "login"
	orderNr := "orderNr"
	sum := float32(42.42)
	params := dto.WithdrawParams{
		Login:   login,
		OrderNr: orderNr,
		Sum:     sum,
	}

	s.withdrawService.err = &appError.InternalError{
		InnerError: fmt.Errorf("database error"),
		Msg:        "service temporary unavailable",
		Code:       appError.ServiceUnavailable,
	}

	apperr := s.action.Withdraw(params)
	s.NotNil(apperr)
	s.Equal(appError.ServiceUnavailable, apperr.Code)

	s.Equal(login, s.withdrawService.param.Login)
	s.Equal(orderNr, s.withdrawService.param.OrderNr)
	s.InDelta(sum, s.withdrawService.param.Sum, 0.001)
}

type orderService struct {
	err    *appError.InternalError
	login  string
	orders []dto.OrderListItem
}

func (t *orderService) GetUserOrders(login string) ([]dto.OrderListItem, *appError.InternalError) {
	t.login = login
	return t.orders, t.err
}

type withdrawService struct {
	param dto.WithdrawParams
	err   *appError.InternalError
}

func (t *withdrawService) Withdraw(params dto.WithdrawParams) *appError.InternalError {
	t.param = params
	return t.err
}
