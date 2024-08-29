package add

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type AddOrderTestSuite struct {
	suite.Suite
	action  *Action
	storage *storage
	logger  *zap.Logger
	logs    *bytes.Buffer
}

func Test_AddOrder(t *testing.T) {
	suite.Run(t, &AddOrderTestSuite{})
}

func (s *AddOrderTestSuite) SetupTest() {
	s.logger, s.logs = logger.SpyLogger()
	s.storage = &storage{}
	s.action = New(s.storage, s.logger)
}

func (s *AddOrderTestSuite) Test_UnableToFindOrderNr() {
	s.storage.findOrderErr = fmt.Errorf("database error")
	err := s.action.SaveOrder("login", "orderNr")
	s.Error(err)
	s.Equal("orderNr", s.storage.orderNrToFind)
}

func (s *AddOrderTestSuite) Test_OrderBelongToDifferentUser() {
	orderNr := "orderNr"

	s.storage.foundOrder = dto.Order{
		OrderNr: orderNr,
		Login:   "different login",
	}

	err := s.action.SaveOrder("login", orderNr)

	s.Error(err)
	s.Equal(err.Code, appError.BadOrderOwnership)
	s.Equal(orderNr, s.storage.orderNrToFind)
}

func (s *AddOrderTestSuite) Test_UserAlreadyRegisteredThisOrder() {
	orderNr := "orderNr"

	s.storage.foundOrder = dto.Order{
		OrderNr: orderNr,
		Login:   "login",
	}

	err := s.action.SaveOrder("login", orderNr)

	s.Error(err)
	s.Equal(err.Code, appError.OrderNrExists)
	s.Equal(orderNr, s.storage.orderNrToFind)
}

func (s *AddOrderTestSuite) Test_UnableToSaveNewOrder() {
	orderNr := "orderNr"

	s.storage.findOrderErr = &order.NotFoundErr{OrderNr: orderNr}
	s.storage.addOrderErr = fmt.Errorf("database error")

	err := s.action.SaveOrder("login", orderNr)

	s.Error(err)
	s.Equal(orderNr, s.storage.orderNrToFind)
	s.Equal(orderNr, s.storage.orderNrToSave)
	s.Equal("login", s.storage.loginToSaveOrderFor)

	s.Contains(s.logs.String(), "database error")
}

func (s *AddOrderTestSuite) Test_OrderSavedSuccessfully() {
	orderNr := "orderNr"

	s.storage.findOrderErr = &order.NotFoundErr{OrderNr: orderNr}

	err := s.action.SaveOrder("login", orderNr)

	s.Nil(err)
	s.Equal(orderNr, s.storage.orderNrToFind)
	s.Equal(orderNr, s.storage.orderNrToSave)
	s.Equal("login", s.storage.loginToSaveOrderFor)
}

type storage struct {
	orderNrToFind string
	findOrderErr  error
	foundOrder    dto.Order

	orderNrToSave       string
	loginToSaveOrderFor string
	addOrderErr         error
}

func (s *storage) FindOrder(orderNr string) (dto.Order, error) {
	s.orderNrToFind = orderNr
	return s.foundOrder, s.findOrderErr
}

func (s *storage) AddOrder(login, orderNr string) error {
	s.loginToSaveOrderFor = login
	s.orderNrToSave = orderNr
	return s.addOrderErr
}
