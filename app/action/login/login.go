package login

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type UserNotFoundErr struct {
	Login string
}

func (e *UserNotFoundErr) Error() string {
	return fmt.Sprintf("user %s not found", e.Login)
}

type Storage interface {
	FindUser(login string) (dto.User, error)
	AddToken(login, token string) error
}

type TokenGenerator interface {
	Generate() string
}

type Login struct {
	storage           Storage
	logger            *zap.Logger
	tokenGenerator    TokenGenerator
	CheckPasswordHash func(password, hash string) bool
}

func New(storage Storage, logger *zap.Logger, tokenGenerator TokenGenerator) *Login {
	return &Login{
		storage:        storage,
		logger:         logger,
		tokenGenerator: tokenGenerator,
	}
}

func (t *Login) Login(params dto.LoginParams) (string, *appError.InternalError) {
	user, err := t.storage.FindUser(params.Login)
	if err != nil {
		var userNotfoundErr *UserNotFoundErr
		if errors.As(err, &userNotfoundErr) {
			return "", &appError.InternalError{
				InnerError: err,
				Msg:        "invalid login/password",
				Code:       appError.InvalidCredentials,
			}
		}

		t.logger.Error(err.Error(), zap.Error(err))

		return "", &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       6,
		}
	}

	if !t.CheckPasswordHash(params.Password, user.PasswordHash) {
		return "", &appError.InternalError{
			InnerError: nil,
			Msg:        "invalid login/password",
			Code:       appError.InvalidCredentials,
		}
	}

	token := t.tokenGenerator.Generate()

	err = t.storage.AddToken(user.Login, token)
	if err != nil {
		t.logger.Error(err.Error(), zap.Error(err))
		return "", &appError.InternalError{
			InnerError: err,
			Msg:        "service temporary unavailable",
			Code:       appError.ServiceUnavailable,
		}
	}

	return token, nil
}
