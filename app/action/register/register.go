package register

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type Storage interface {
	IsLoginExists(login string) (bool, error)
	SaveUser(login, password string) error
	AddToken(login, token string) error
}

type TokenGenerator interface {
	Generate() string
}

type Register struct {
	storage        Storage
	logger         *zap.Logger
	tokenGenerator TokenGenerator
}

func New(storage Storage, logger *zap.Logger, tokenGenerator TokenGenerator) *Register {
	return &Register{
		storage:        storage,
		logger:         logger,
		tokenGenerator: tokenGenerator,
	}
}

func (r *Register) Register(params dto.RegisterParams) (string, *appError.InternalError) {
	isExists, err := r.storage.IsLoginExists(params.Login)
	if err != nil {
		err := fmt.Errorf("unable to check loginToSave exists: %w", err)
		r.logger.Error(err.Error(), zap.Error(err))

		return "", &appError.InternalError{
			InnerError: err,
			Msg:        "unable to check loginToSave exists",
			Code:       appError.UnableToCheckLoginExists,
		}
	}

	if isExists {
		return "", &appError.InternalError{
			InnerError: err,
			Msg:        "login already in use",
			Code:       appError.LoginAlreadyInUse,
		}
	}

	err = r.storage.SaveUser(params.Login, params.Password)
	if err != nil {
		msg := "unable to save user"
		r.logger.Error(msg, zap.Error(err))

		return "", &appError.InternalError{
			InnerError: err,
			Msg:        msg,
			Code:       appError.UnableToSaveUser,
		}
	}

	token := r.tokenGenerator.Generate()

	err = r.storage.AddToken(params.Login, token)
	if err != nil {
		msg := "unable to save token"
		r.logger.Error(msg, zap.Error(err))

		return "", &appError.InternalError{
			InnerError: err,
			Msg:        msg,
			Code:       appError.UnableToSaveToken,
		}
	}

	return token, nil
}
