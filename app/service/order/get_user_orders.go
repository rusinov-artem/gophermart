package order

import (
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

func (t *Service) GetUserOrders(login string) ([]dto.OrderListItem, *appError.InternalError) {
	logger := t.logger.With(zap.String("login", login))

	orders, err := t.storage.ListOrders(login)
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return nil, &appError.InternalError{
			InnerError: err,
			Msg:        "service unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	if len(orders) == 0 {
		return nil, &appError.InternalError{
			Msg:  "there are no orders",
			Code: appError.NoOrdersFound,
		}
	}

	err = t.accrual.EnrichOrders(orders)
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return nil, &appError.InternalError{
			InnerError: err,
			Msg:        "service unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	return orders, nil
}
