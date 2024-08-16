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

type GetWithdrawalsAction interface {
	GetWithdrawals(login string) ([]dto.Withdrawal, *appError.InternalError)
}

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	h.auth(w, r, func(ctx context.Context, login string) {
		w.Header().Set("Content-Type", "application/json")
		withdrawals, internalErr := h.GetWithdrawalsAction(ctx).GetWithdrawals(login)
		if internalErr != nil {
			converter.ConvertError(w, internalErr)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(encodeWithdrawals(withdrawals))
	})
}

func encodeWithdrawals(withdrawals []dto.Withdrawal) []byte {
	type jsonWithdrawals struct {
		Order       string  `json:"order"`
		Sum         float32 `json:"sum"`
		ProcessedAt string  `json:"processed_at"`
	}

	data := make([]jsonWithdrawals, 0, len(withdrawals))
	for i := range withdrawals {
		data = append(data, jsonWithdrawals{
			Order: withdrawals[i].OrderNr,
			Sum:   withdrawals[i].Sum,
			ProcessedAt: func() string {
				return withdrawals[i].ProcessedAt.Format(time.RFC3339)
			}(),
		})
	}

	b, _ := json.Marshal(data)
	return b
}
