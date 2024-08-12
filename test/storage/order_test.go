package storage

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/storage"
	"github.com/rusinov-artem/gophermart/test"
)

type OrderStorageTestSuite struct {
	suite.Suite
	ctx     context.Context
	pool    *pgxpool.Pool
	storage *storage.RegistrationStorage
}

func Test_OrderStorate(t *testing.T) {
	suite.Run(t, &OrderStorageTestSuite{})
}

func (s *OrderStorageTestSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()
	dsn := test.CreateTestDB("test_order_storage")
	s.pool, err = pgxpool.New(context.Background(), dsn)
	s.Require().NoError(err)
}

func (s *OrderStorageTestSuite) SetupTest() {
	s.T().Parallel()
	s.storage = storage.NewRegistrationStorage(s.ctx, s.pool)
}

func (s *OrderStorageTestSuite) Test_OrderNotFound() {
	_, err := s.storage.FindOrder("unknownOrderNr")

	var notFoundErr *order.NotFoundErr
	s.ErrorAs(err, &notFoundErr)
}

func (s *OrderStorageTestSuite) Test_CanAddOrder() {
	s.Require().NoError(s.storage.SaveUser("login", "password"))
	s.Require().NoError(s.storage.AddOrder("login", "OrderNR001"))

	dtoOrder, err := s.storage.FindOrder("OrderNR001")
	s.Require().NoError(err)

	s.Equal("login", dtoOrder.Login)
	s.Equal("OrderNR001", dtoOrder.OrderNr)
}

func (s *OrderStorageTestSuite) Test_CantAddSameOrder() {
	s.Require().NoError(s.storage.SaveUser("login2", "password"))
	s.Require().NoError(s.storage.AddOrder("login2", "OrderNR002"))

	_, err := s.storage.FindOrder("OrderNR002")
	s.Require().NoError(err)

	s.Require().Error(s.storage.AddOrder("login2", "OrderNR002"))
}

func (s *OrderStorageTestSuite) Test_UserCanHaveMultipleOrders() {
	s.Require().NoError(s.storage.SaveUser("login3", "password"))
	s.Require().NoError(s.storage.AddOrder("login3", "OrderNR003"))
	s.Require().NoError(s.storage.AddOrder("login3", "OrderNR004"))

	_, err := s.storage.FindOrder("OrderNR004")
	s.Require().NoError(err)

	_, err = s.storage.FindOrder("OrderNR003")
	s.Require().NoError(err)
}

func (s *OrderStorageTestSuite) Test_CantAddSingleOrderToMultipleUsers() {
	s.Require().NoError(s.storage.SaveUser("login4", "password"))
	s.Require().NoError(s.storage.SaveUser("login5", "password"))

	s.Require().NoError(s.storage.AddOrder("login4", "OrderNR006"))
	s.Require().Error(s.storage.AddOrder("login5", "OrderNR006"))
}
