package accrual

import (
	"fmt"

	"go.uber.org/zap"

	appOrder "github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
)

const REGISTERED = "REGISTERED"

type Client interface {
	GetSingleOrder(orderNr string) (dto.OrderListItem, error)
}

type Storage interface {
	UpdateOrdersState(orders []dto.OrderListItem) error
}

type Service struct {
	client  Client
	logger  *zap.Logger
	storage Storage
}

func NewService(client Client, storage Storage, logger *zap.Logger) *Service {
	return &Service{
		client:  client,
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) EnrichOrders(orders []dto.OrderListItem) error {
	for i := range orders {
		s.enrichSingleOrder(&orders[i])
	}

	err := s.storage.UpdateOrdersState(orders)
	if err != nil {
		err := fmt.Errorf("unable to enrich orders: %w", err)
		s.logger.Error(err.Error(), zap.Error(err))

		return err
	}

	return nil
}

func (s *Service) enrichSingleOrder(order *dto.OrderListItem) {
	accrualOrder, err := s.client.GetSingleOrder(order.OrderNr)
	if err != nil {
		s.logger.Error(err.Error(), zap.Error(err))
		return
	}

	order.Accrual = accrualOrder.Accrual
	order.Status = accrualOrder.Status
	if order.Status == "REGISTERED" {
		order.Status = appOrder.NEW
	}
}
