package order

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type OrderServiceTestSuite struct {
	suite.Suite
	logger  *zap.Logger
	logs    *bytes.Buffer
	accrual *accrual
	storage *storage
	action  *Service
}

func Test_OrderService(t *testing.T) {
	suite.Run(t, &OrderServiceTestSuite{})
}

func (s *OrderServiceTestSuite) SetupTest() {
	s.logger, s.logs = logger.SpyLogger()
	s.accrual = &accrual{}
	s.storage = &storage{}
	s.action = NewOrderService(s.logger, s.storage, s.accrual)
}

func (s *OrderServiceTestSuite) Test_UnableToFetchOrdersFromDB() {
	login := "login"
	s.storage.GetUserOrdersErr = fmt.Errorf("db error")

	_, err := s.action.GetUserOrders(login)
	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(login, s.storage.login)

	s.Contains(s.logs.String(), "db error")
}

func (s *OrderServiceTestSuite) Test_UserDoesNotHaveOrders() {
	login := "login"

	orders, err := s.action.GetUserOrders(login)
	s.NotNil(err)
	s.Equal(appError.NoOrdersFound, err.Code)
	s.Equal(0, len(orders))
	s.Equal(login, s.storage.login)
}

func (s *OrderServiceTestSuite) Test_UnableToFetchAccrual() {
	login := "login"

	s.storage.foundOrders = []dto.OrderListItem{
		{
			OrderNr:  "11111",
			Status:   "NEW",
			UploadAt: time.Time{},
		},
		{
			OrderNr:  "22222",
			Status:   "PROCESSED",
			Accrual:  300,
			UploadAt: time.Time{},
		},
	}

	s.accrual.fetchOrdersErr = fmt.Errorf("accrual error")

	_, err := s.action.GetUserOrders(login)

	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(login, s.storage.login)
	s.Equal(s.storage.foundOrders, s.accrual.queryOrders)
	s.Contains(s.logs.String(), "accrual error")
}

func (s *OrderServiceTestSuite) Test_Success() {
	login := "login"

	s.storage.foundOrders = []dto.OrderListItem{
		{
			OrderNr:  "11111",
			Status:   "NEW",
			UploadAt: time.Time{},
		},
		{
			OrderNr:  "22222",
			Status:   "PROCESSED",
			Accrual:  300,
			UploadAt: time.Time{},
		},
	}

	_, err := s.action.GetUserOrders(login)

	s.Nil(err)
	s.Equal(login, s.storage.login)
	s.Equal(s.storage.foundOrders, s.accrual.queryOrders)
}

type storage struct {
	login            string
	GetUserOrdersErr error
	foundOrders      []dto.OrderListItem
}

func (s *storage) ListOrders(login string) ([]dto.OrderListItem, error) {
	s.login = login
	return s.foundOrders, s.GetUserOrdersErr
}

type accrual struct {
	fetchOrdersErr error
	queryOrders    []dto.OrderListItem
}

func (t *accrual) EnrichOrders(orders []dto.OrderListItem) error {
	t.queryOrders = orders
	return t.fetchOrdersErr
}
