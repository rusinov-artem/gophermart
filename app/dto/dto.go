package dto

import (
	"time"

	"github.com/rusinov-artem/gophermart/app/action/order"
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

type OrderList []OrderListItem

func (t OrderList) Total() float32 {
	var total float32
	for i := range t {
		if t[i].Status == order.PROCESSED {
			total += t[i].Accrual
		}
	}

	return total
}

type Balance struct {
	Current   float32
	Withdrawn float32
}
