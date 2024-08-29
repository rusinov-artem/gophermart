package get

import (
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type BalanceStorage interface {
	Balance(login string) (float32, error)
	Withdrawn(login string) (float32, error)
}

type OrderService interface {
	GetUserOrders(login string) ([]dto.OrderListItem, *appError.InternalError)
}

type Action struct {
	service OrderService
	storage BalanceStorage
	logger  *zap.Logger
}

func New(service OrderService, storage BalanceStorage, logger *zap.Logger) *Action {
	return &Action{
		service: service,
		storage: storage,
		logger:  logger,
	}
}

func (t *Action) GetBalance(login string) (dto.Balance, *appError.InternalError) {
	logger := t.logger.With(zap.String("login", login))

	result := dto.Balance{}
	_, appErr := t.service.GetUserOrders(login)
	if appErr != nil {
		return result, appErr
	}

	current, err := t.storage.Balance(login)
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return result, &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}
	result.Current = current

	withdrawn, err := t.storage.Withdrawn(login)
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return result, &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}
	result.Withdrawn = withdrawn

	return result, nil
}
