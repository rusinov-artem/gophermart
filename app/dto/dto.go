package dto

import "time"

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

type Order struct {
	OrderNr string
	Login   string
}

type OrderListItem struct {
	OrderNr  string
	Status   string
	Accrual  int64
	UploadAt time.Time
}
