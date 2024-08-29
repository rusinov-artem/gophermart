package withdraw

import (
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type Transaction interface {
	Begin() error
	Rollback()
	LockUser() error
	AvailablePoints() (float32, error)
	Withdraw(orderNr string, amount float32) error
	Commit() error
}

type TransactionFactory func(login string) Transaction

type Service struct {
	newTx  TransactionFactory
	logger *zap.Logger
}

func NewWithdrawService(txFactory TransactionFactory, logger *zap.Logger) *Service {
	return &Service{
		newTx:  txFactory,
		logger: logger,
	}
}

func (s *Service) Withdraw(params dto.WithdrawParams) *appError.InternalError {
	logger := s.logger.With(
		zap.String("login", params.Login),
	)

	tx := s.newTx(params.Login)
	err := tx.Begin()
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	defer tx.Rollback()

	err = tx.LockUser()
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	availablePoints, err := tx.AvailablePoints()
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	if availablePoints < params.Sum {
		return &appError.InternalError{
			InnerError: err,
			Msg:        "not enough points",
			Code:       appError.NotEnoughPoints,
		}
	}

	err = tx.Withdraw(params.OrderNr, params.Sum)
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error(), zap.Error(err))
		return &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	return nil
}
