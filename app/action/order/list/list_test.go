package list

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

type ListOrderActionTestSuite struct {
	suite.Suite
	action  *Action
	storage *storage
	logger  *zap.Logger
	logs    *bytes.Buffer
	accrual *accrual
}

func Test_ListOrderAction(t *testing.T) {
	suite.Run(t, &ListOrderActionTestSuite{})
}

func (s *ListOrderActionTestSuite) SetupTest() {
	s.logger, s.logs = logger.SpyLogger()
	s.accrual = &accrual{}
	s.storage = &storage{}
	s.action = New(s.storage, s.accrual, s.logger)
}

func (s *ListOrderActionTestSuite) Test_UnableToFetchOrdersFromDB() {
	login := "login"
	s.storage.listOrdersErr = fmt.Errorf("db error")

	_, err := s.action.ListOrders(login)
	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(login, s.storage.login)

	s.Contains(s.logs.String(), "db error")
}

func (s *ListOrderActionTestSuite) Test_UserDoesNotHaveOrders() {
	login := "login"

	orders, err := s.action.ListOrders(login)
	s.NotNil(err)
	s.Equal(appError.NoOrdersFound, err.Code)
	s.Equal(0, len(orders))
	s.Equal(login, s.storage.login)
}

func (s *ListOrderActionTestSuite) Test_UnableToFetchAccrual() {
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

	_, err := s.action.ListOrders(login)

	s.NotNil(err)
	s.Equal(appError.ServiceUnavailable, err.Code)
	s.Equal(login, s.storage.login)
	s.Equal(s.storage.foundOrders, s.accrual.queryOrders)
	s.Contains(s.logs.String(), "accrual error")
}

func (s *ListOrderActionTestSuite) Test_Success() {
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

	_, err := s.action.ListOrders(login)

	s.Nil(err)
	s.Equal(login, s.storage.login)
	s.Equal(s.storage.foundOrders, s.accrual.queryOrders)
}

type storage struct {
	login         string
	listOrdersErr error
	foundOrders   []dto.OrderListItem
}

func (s *storage) ListOrders(login string) ([]dto.OrderListItem, error) {
	s.login = login
	return s.foundOrders, s.listOrdersErr
}

type accrual struct {
	fetchOrdersErr error
	queryOrders    []dto.OrderListItem
}

func (t *accrual) FetchOrders(orders *[]dto.OrderListItem) error {
	t.queryOrders = *orders
	return t.fetchOrdersErr
}
