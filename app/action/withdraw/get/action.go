package get

import (
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type WithdrawalsStorage interface {
	GetWithdrawals(login string) ([]dto.Withdrawal, error)
}

type Action struct {
	logger  *zap.Logger
	storage WithdrawalsStorage
}

func New(logger *zap.Logger, storage WithdrawalsStorage) *Action {
	return &Action{
		logger:  logger,
		storage: storage,
	}
}

func (a *Action) GetWithdrawals(login string) ([]dto.Withdrawal, *appError.InternalError) {
	withdrawals, err := a.storage.GetWithdrawals(login)
	if err != nil {
		return nil, &appError.InternalError{
			InnerError: err,
			Msg:        "Service temporarily unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	if len(withdrawals) == 0 {
		return nil, &appError.InternalError{
			Msg:  "No withdrawals found",
			Code: appError.NoWithdrawals,
		}
	}

	return withdrawals, nil
}
