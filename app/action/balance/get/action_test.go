package get

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type GetBalanceActionTestSuite struct {
	suite.Suite
	action       *Action
	orderService *orderService
}

func Test_GetBalanceAction(t *testing.T) {
	suite.Run(t, &GetBalanceActionTestSuite{})
}

func (s *GetBalanceActionTestSuite) SetupTest() {
	s.orderService = &orderService{}
	s.action = New(s.orderService)
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

func (s *GetBalanceActionTestSuite) Test_CanCalculateBalance() {
	login := "login"
	s.orderService.orders = []dto.OrderListItem{
		{
			Status:  order.NEW,
			Accrual: 987,
		},
		{
			Status:  order.PROCESSED,
			Accrual: 42.42,
		},
	}

	balance, appErr := s.action.GetBalance(login)
	s.Nil(appErr)
	s.Equal(login, s.orderService.ordersForLogin)
	s.Equal(s.orderService.err, appErr)
	s.InDelta(float32(42.42), balance.Current, 0.001)
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
