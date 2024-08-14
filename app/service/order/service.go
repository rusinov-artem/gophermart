package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
)

type Accrual interface {
	EnrichOrders(orders []dto.OrderListItem) error
}

type Storage interface {
	ListOrders(login string) ([]dto.OrderListItem, error)
}

type Service struct {
	ctx     context.Context
	logger  *zap.Logger
	storage Storage
	accrual Accrual
}

func NewOrderService(logger *zap.Logger, storage Storage, accrual Accrual) *Service {
	return &Service{
		logger:  logger,
		storage: storage,
		accrual: accrual,
	}
}
