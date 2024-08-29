package handler

import (
	"context"
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
	h.auth(w, r, func(ctx context.Context, login string) {
		w.Header().Set("Content-Type", "application/json")
		orders, internalErr := h.ListOrdersAction(ctx).ListOrders(login)
		if internalErr != nil {
			converter.ConvertError(w, internalErr)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(encodeOrderList(orders))
	})
}

func encodeOrderList(orders []dto.OrderListItem) []byte {
	type jsonOrder struct {
		Number     string   `json:"number"`
		Status     string   `json:"status"`
		Accrual    *float32 `json:"accrual,omitempty"`
		UploadedAt string   `json:"uploaded_at"`
	}

	jsonOrderList := make([]jsonOrder, len(orders))
	for i := range orders {
		jsonOrderList[i] = jsonOrder{
			Number: orders[i].OrderNr,
			Status: orders[i].Status,

			Accrual: func() *float32 {
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
