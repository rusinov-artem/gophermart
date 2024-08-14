package list

import (
	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type OrderService interface {
	GetUserOrders(login string) ([]dto.OrderListItem, *appError.InternalError)
}

type Action struct {
	service OrderService
}

func New(service OrderService) *Action {
	return &Action{
		service: service,
	}
}

func (t *Action) ListOrders(login string) ([]dto.OrderListItem, *appError.InternalError) {
	return t.service.GetUserOrders(login)
}
