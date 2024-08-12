package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/app/http/converter"
)

type ListOrdersAction interface {
	ListOrders(login string) ([]dto.OrderListItem, *appError.InternalError)
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.Context(r.Context())
	defer cancelFN()

	w.Header().Set("Content-Type", "application/json")

	auth := h.AuthService(ctx)
	login, err := auth.Auth(getToken(r))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	action := h.ListOrdersAction(ctx)
	orders, internalErr := action.ListOrders(login)
	if internalErr != nil {
		converter.ConvertError(w, internalErr)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(encodeOrderList(orders))
}

func encodeOrderList(orders []dto.OrderListItem) []byte {
	type jsonOrder struct {
		Number     string `json:"number"`
		Status     string `json:"status"`
		Accrual    *int64 `json:"accrual,omitempty"`
		UploadedAt string `json:"uploaded_at"`
	}

	jsonOrderList := make([]jsonOrder, len(orders))
	for i := range orders {
		jsonOrderList[i] = jsonOrder{
			Number: orders[i].OrderNr,
			Status: orders[i].Status,

			Accrual: func() *int64 {
				if orders[i].Status == "PROCESSED" {
					return &orders[i].Accrual
				}

				return nil
			}(),

			UploadedAt: func() string {
				return orders[i].UploadAt.Format(time.RFC3339)
			}(),
		}
	}

	b, _ := json.MarshalIndent(jsonOrderList, "", "\t")
	return b
}
