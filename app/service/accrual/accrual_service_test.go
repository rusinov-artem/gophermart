package accrual

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	appOrder "github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type AccrualServiceTestSuite struct {
	suite.Suite
	service *Service
	client  *client
	logger  *zap.Logger
	logs    *bytes.Buffer
	storage *storage
}

func Test_AccrualService(t *testing.T) {
	suite.Run(t, &AccrualServiceTestSuite{})
}

func (s *AccrualServiceTestSuite) SetupTest() {
	s.logger, s.logs = logger.SpyLogger()
	s.client = &client{}
	s.storage = &storage{}
	s.service = NewService(s.client, s.storage, s.logger)
}

func (s *AccrualServiceTestSuite) Test_FetchEmptyList() {
	s.Require().NoError(s.service.EnrichOrders(nil))
}

func (s *AccrualServiceTestSuite) Test_UnableGetOrderFromAccrual() {
	dt := time.Now()
	orders := []dto.OrderListItem{
		{
			OrderNr:  "OrderNr001",
			Status:   appOrder.NEW,
			Accrual:  0,
			UploadAt: dt,
		},
	}

	s.client.err = fmt.Errorf("network error")

	err := s.service.EnrichOrders(orders)
	s.Require().NoError(err)

	s.Equal(dto.OrderListItem{
		OrderNr:  "OrderNr001",
		Status:   appOrder.NEW,
		Accrual:  0,
		UploadAt: dt,
	}, orders[0])

	s.Contains(s.logs.String(), "network error")
}

func (s *AccrualServiceTestSuite) Test_UnableGetOrderFromAccrual_OrderStateNotChange() {
	dt := time.Now()
	orders := []dto.OrderListItem{
		{
			OrderNr:  "OrderNr001",
			Status:   appOrder.PROCESSED,
			Accrual:  55,
			UploadAt: dt,
		},
	}

	s.client.err = fmt.Errorf("network error")

	err := s.service.EnrichOrders(orders)
	s.Require().NoError(err)

	s.Equal(dto.OrderListItem{
		OrderNr:  "OrderNr001",
		Status:   appOrder.PROCESSED,
		Accrual:  55,
		UploadAt: dt,
	}, orders[0])

	s.Contains(s.logs.String(), "network error")
}

func (s *AccrualServiceTestSuite) Test_DoNotChangeStatusToREGISTERED() {
	dt := time.Now()
	orders := []dto.OrderListItem{
		{
			OrderNr:  "OrderNr001",
			Status:   appOrder.NEW,
			Accrual:  55,
			UploadAt: dt,
		},
	}

	s.client.order = dto.OrderListItem{
		OrderNr:  "OrderNr001",
		Status:   REGISTERED,
		Accrual:  55,
		UploadAt: dt,
	}

	err := s.service.EnrichOrders(orders)
	s.Require().NoError(err)

	s.Equal(dto.OrderListItem{
		OrderNr:  "OrderNr001",
		Status:   appOrder.NEW,
		Accrual:  55,
		UploadAt: dt,
	}, orders[0])

}

func (s *AccrualServiceTestSuite) Test_UnableToUpdateOrderState() {
	dt := time.Now()
	orders := []dto.OrderListItem{
		{
			OrderNr:  "OrderNr001",
			Status:   "NEW",
			Accrual:  0,
			UploadAt: dt,
		},
	}

	s.storage.err = fmt.Errorf("database error")
	err := s.service.EnrichOrders(orders)

	s.Error(err)
	s.Contains(s.logs.String(), "database error")
}

func (s *AccrualServiceTestSuite) Test_SingleOrder() {
	dt := time.Now()
	orders := []dto.OrderListItem{
		{
			OrderNr:  "OrderNr001",
			Status:   "NEW",
			Accrual:  0,
			UploadAt: dt,
		},
	}

	s.client.order = dto.OrderListItem{
		OrderNr: "OrderNr001",
		Status:  "PROCESSED",
		Accrual: 42,
	}

	err := s.service.EnrichOrders(orders)
	s.Require().NoError(err)

	s.Equal(dto.OrderListItem{
		OrderNr:  "OrderNr001",
		Status:   "PROCESSED",
		Accrual:  42,
		UploadAt: dt,
	}, orders[0])
}

type client struct {
	orderNr string
	err     error
	order   dto.OrderListItem
}

func (c *client) GetSingleOrder(orderNr string) (dto.OrderListItem, error) {
	c.orderNr = orderNr
	return c.order, c.err
}

type storage struct {
	ordersToUpdate []dto.OrderListItem
	err            error
}

func (s *storage) UpdateOrdersState(orders []dto.OrderListItem) error {
	s.ordersToUpdate = orders
	return s.err
}
