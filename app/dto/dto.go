package dto

import (
	"time"
)

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
	Accrual  float32
	UploadAt time.Time
}

type Balance struct {
	Current   float32
	Withdrawn float32
}

type WithdrawParams struct {
	Login   string
	OrderNr string
	Sum     float32
}

type Withdrawal struct {
	OrderNr     string
	Sum         float32
	ProcessedAt time.Time
}
