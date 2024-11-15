package add

import (
	"errors"

	"go.uber.org/zap"

	actionError "github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type Storage interface {
	FindOrder(orderNr string) (dto.Order, error)
	AddOrder(login, orderNr string) error
}

type Action struct {
	storage Storage
	logger  *zap.Logger
}

func New(storage Storage, logger *zap.Logger) *Action {
	return &Action{

		storage: storage,
		logger:  logger,
	}
}

func (t *Action) SaveOrder(login, orderNr string) *appError.InternalError {
	logger := t.logger.With(
		zap.String("login", login),
		zap.String("orderNr", orderNr),
	)

	order, err := t.storage.FindOrder(orderNr)
	if err != nil {
		var orderNotFoundErr *actionError.NotFoundErr
		if errors.As(err, &orderNotFoundErr) {
			return t.addOrder(logger, login, orderNr)
		}

		logger.Error(err.Error(), zap.Error(err))

		return &appError.InternalError{
			InnerError: err,
			Msg:        "service unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	if order.Login != login {
		return &appError.InternalError{
			Msg:  "order belongs to different login",
			Code: appError.BadOrderOwnership,
		}
	}

	return &appError.InternalError{
		Msg:  "order already exists",
		Code: appError.OrderNrExists,
	}
}

func (t *Action) addOrder(logger *zap.Logger, login, orderNr string) *appError.InternalError {
	err := t.storage.AddOrder(login, orderNr)
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))

		return &appError.InternalError{
			InnerError: err,
			Msg:        "service unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	return nil
}
