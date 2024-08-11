package register

import (
	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

func (r *Register) Validate(params dto.RegisterParams) *appError.ValidationError {
	fields := make(map[string][]string)

	hasErrors := false
	loginErrors := r.validateLogin(params.Login)
	if len(loginErrors) > 0 {
		hasErrors = true
		fields["loginToSave"] = loginErrors
	}

	passwordErrors := r.validatePassword(params.Password)
	if len(passwordErrors) > 0 {
		hasErrors = true
		fields["password"] = passwordErrors
	}

	if hasErrors {
		return &appError.ValidationError{Fields: fields}
	}

	return nil
}

func (r *Register) validateLogin(login string) []string {
	var errors []string

	if login == "" {
		errors = append(errors, "loginToSave is required")
		return errors
	}

	return errors
}

func (r *Register) validatePassword(password string) []string {
	var errors []string

	if password == "" {
		errors = append(errors, "password is required")
		return errors
	}

	if len([]rune(password)) < 6 {
		errors = append(errors, "password is too short")
		return errors
	}

	return errors
}
