package storage

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
	"github.com/rusinov-artem/gophermart/app/storage"
	"github.com/rusinov-artem/gophermart/test"
)

type OrderStorageTestSuite struct {
	suite.Suite
	ctx     context.Context
	pool    *pgxpool.Pool
	storage *storage.Storage
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
	s.storage = storage.NewStorage(s.ctx, s.pool)
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

func (s *OrderStorageTestSuite) Test_CanGetEmptyOrdersList() {
	orders, err := s.storage.ListOrders("login6")
	s.Require().NoError(err)
	s.Len(orders, 0)
}

func (s *OrderStorageTestSuite) Test_CanListSingleOrder() {
	s.Require().NoError(s.storage.SaveUser("login7", "password"))
	s.Require().NoError(s.storage.AddOrder("login7", "OrderNR007"))
	orders, err := s.storage.ListOrders("login7")
	s.Require().NoError(err)
	s.Len(orders, 1)
	s.Equal("OrderNR007", orders[0].OrderNr)
	s.Equal("NEW", orders[0].Status)
	s.Equal(float32(0), orders[0].Accrual)
	s.NotEmpty(orders[0].UploadAt)
}

func (s *OrderStorageTestSuite) Test_CanListMultipleOrders() {
	s.Require().NoError(s.storage.SaveUser("login8", "password"))
	s.Require().NoError(s.storage.AddOrder("login8", "OrderNR008"))
	s.Require().NoError(s.storage.AddOrder("login8", "OrderNR009"))
	s.Require().NoError(s.storage.AddOrder("login8", "OrderNR010"))

	s.SetOrderUploadAt("OrderNR008", "2024-04-01 00:00:00")
	s.SetOrderUploadAt("OrderNR009", "2024-05-01 00:00:00")
	s.SetOrderUploadAt("OrderNR010", "2024-05-06 00:00:00")

	orders, err := s.storage.ListOrders("login8")

	s.Require().NoError(err)
	s.Len(orders, 3)

	s.LessOrEqual(orders[0].UploadAt, orders[1].UploadAt)
}

func (s *OrderStorageTestSuite) Test_CanUpdateSingleOrderState() {
	s.Require().NoError(s.storage.SaveUser("login9", "password"))
	s.Require().NoError(s.storage.AddOrder("login9", "OrderNR011"))

	err := s.storage.UpdateOrdersState([]dto.OrderListItem{
		{
			OrderNr: "OrderNR011",
			Status:  "PROCESSED",
			Accrual: 42,
		},
	})
	s.Require().NoError(err)

	orders, err := s.storage.ListOrders("login9")
	s.Require().NoError(err)
	s.Len(orders, 1)
	s.Equal("OrderNR011", orders[0].OrderNr)
	s.Equal("PROCESSED", orders[0].Status)
	s.Equal(float32(42), orders[0].Accrual)
}

func (s *OrderStorageTestSuite) Test_CanUpdateMultipleOrders() {
	s.Require().NoError(s.storage.SaveUser("login10", "password"))

	s.Require().NoError(s.storage.AddOrder("login10", "OrderNR012"))
	s.SetOrderUploadAt("OrderNR012", "2024-04-01 00:00:00")

	s.Require().NoError(s.storage.AddOrder("login10", "OrderNR013"))
	s.SetOrderUploadAt("OrderNR012", "2024-04-02 00:00:00")

	err := s.storage.UpdateOrdersState([]dto.OrderListItem{
		{
			OrderNr: "OrderNR012",
			Status:  "PROCESSED",
			Accrual: 42,
		},
		{
			OrderNr: "OrderNR013",
			Status:  "PROCESSING",
		},
	})

	s.Require().NoError(err)

	orders, err := s.storage.ListOrders("login10")
	s.Require().NoError(err)
	s.Len(orders, 2)
	s.Equal(float32(42), orders[0].Accrual)
	s.Equal(float32(0), orders[1].Accrual)
}

func (s *OrderStorageTestSuite) SetOrderUploadAt(orderNr string, dt string) {
	sqlStr := `UPDATE "order" SET upload_at = $1 WHERE order_nr = $2`

	_, err := s.pool.Exec(s.ctx, sqlStr, dt, orderNr)
	s.Require().NoError(err)
}
