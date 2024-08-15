package withdraw

import (
	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type OrderService interface {
	GetUserOrders(login string) ([]dto.OrderListItem, *appError.InternalError)
}

type WithdrawService interface {
	Withdraw(params dto.WithdrawParams) *appError.InternalError
}

type Action struct {
	orderService    OrderService
	withdrawService WithdrawService
}

func New(orderService OrderService, withdrawService WithdrawService) *Action {
	return &Action{
		orderService:    orderService,
		withdrawService: withdrawService,
	}
}

func (t *Action) Withdraw(params dto.WithdrawParams) *appError.InternalError {
	err := t.updatePointsInfo(params.Login)
	if err != nil {
		return err
	}

	return t.withdrawService.Withdraw(params)
}

func (t *Action) updatePointsInfo(login string) *appError.InternalError {
	_, internalErr := t.orderService.GetUserOrders(login)
	if internalErr != nil {
		if internalErr.Code == appError.NoOrdersFound {
			return &appError.InternalError{
				Msg:  "not enough points",
				Code: appError.NotEnoughPoints,
			}
		}

		return internalErr
	}

	return nil
}
