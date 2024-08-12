package login

import (
	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

func (t *Login) Validate(params dto.LoginParams) *appError.ValidationError {
	fields := make(map[string][]string)

	hasErrors := false
	loginErrors := t.validateLogin(params.Login)
	if len(loginErrors) > 0 {
		hasErrors = true
		fields["login"] = loginErrors
	}

	passwordErrors := t.validatePassword(params.Password)
	if len(passwordErrors) > 0 {
		hasErrors = true
		fields["password"] = passwordErrors
	}

	if hasErrors {
		return &appError.ValidationError{Fields: fields}
	}

	return nil
}

func (t *Login) validateLogin(login string) []string {
	var errors []string

	if login == "" {
		errors = append(errors, "login is required")
		return errors
	}

	return errors
}

func (t *Login) validatePassword(password string) []string {
	var errors []string

	if password == "" {
		errors = append(errors, "password is required")
		return errors
	}

	return errors
}
