package storage

import (
	"fmt"

	actionErr "github.com/rusinov-artem/gophermart/app/action/login"
	"github.com/rusinov-artem/gophermart/app/dto"
)

func (r *RegistrationStorage) FindUser(login string) (dto.User, error) {
	sqlStr := `
	SELECT login, password_hash FROM "user" WHERE login = $1`

	user := dto.User{}

	row, err := r.pool.Query(r.ctx, sqlStr, login)
	if err != nil {
		return user, fmt.Errorf("unable to find user: %w", err)
	}

	if !row.Next() {
		return user, &actionErr.UserNotFoundErr{Login: login}
	}

	err = row.Scan(&user.Login, &user.PasswordHash)
	if err != nil {
		return user, fmt.Errorf("unable to find user: %w", err)
	}

	return user, nil
}