package dto

type RegisterParams struct {
	Login    string
	Password string
}

type LoginParams struct {
	Login    string
	Password string
}

type User struct {
	Login        string
	PasswordHash string
}
