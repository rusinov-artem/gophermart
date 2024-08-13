package clients

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/dto"
	"github.com/rusinov-artem/gophermart/app/service/accrual"
	accrualClient "github.com/rusinov-artem/gophermart/app/service/accrual/client"
	"github.com/rusinov-artem/gophermart/test/fakes"
)

type ClientTestSuite struct {
	suite.Suite
	client  *accrualClient.Client
	accrual *fakes.Accrual
	ctx     context.Context
}

func Test_Client(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}

func (s *ClientTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.accrual = fakes.NewAccrual(s.T())
	s.client = accrualClient.New(s.ctx, s.accrual.URL())
}

func (s *ClientTestSuite) Test_ErrorIfOrderNotFound() {
	s.accrual.WillReturn204()
	_, err := s.client.GetSingleOrder("unknown")
	s.Error(err)
	s.Equal(s.accrual.Req.Method, http.MethodGet)
	s.Equal("/api/orders/unknown", s.accrual.Req.Path)
}

func (s *ClientTestSuite) Test_CanFetchOrder() {
	orderNr := "OrderNR"
	points := int64(123)
	s.accrual.WillReturnOrder(dto.OrderListItem{
		OrderNr: orderNr,
		Status:  accrual.REGISTERED,
		Accrual: points,
	})

	order, err := s.client.GetSingleOrder("orderNr")
	s.NoError(err)
	s.Equal(s.accrual.Req.Method, http.MethodGet)
	s.Equal("/api/orders/orderNr", s.accrual.Req.Path)

	s.Equal(accrual.REGISTERED, order.Status)
	s.Equal(orderNr, order.OrderNr)
	s.Equal(points, order.Accrual)
}
