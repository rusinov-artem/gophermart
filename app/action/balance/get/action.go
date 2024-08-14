package get

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

func (t *Action) GetBalance(login string) (dto.Balance, *appError.InternalError) {
	balance := dto.Balance{}
	orders, appErr := t.service.GetUserOrders(login)
	if appErr != nil {
		return balance, appErr
	}

	balance.Current = dto.OrderList(orders).Total()
	balance.Withdrawn = 0

	return balance, nil
}
