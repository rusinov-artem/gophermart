package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/app/http/converter"
)

type WithdrawAction interface {
	Withdraw(params dto.WithdrawParams) *appError.InternalError
}

type withdrawJSON struct {
	OrderNr string  `json:"number"`
	Sum     float32 `json:"sum"`
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.Context(r.Context())
	defer cancelFN()
	w.Header().Set("Content-Type", "application/json")

	var envelop withdrawJSON
	err := json.NewDecoder(r.Body).Decode(&envelop)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//err = order.ValidateOrderNr(envelop.OrderNr)
	//if err != nil {
	//	w.WriteHeader(http.StatusUnprocessableEntity)
	//	return
	//}

	auth := h.AuthService(ctx)
	login, err := auth.Auth(getToken(r))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := dto.WithdrawParams{
		Login:   login,
		OrderNr: envelop.OrderNr,
		Sum:     envelop.Sum,
	}

	internalErr := h.WithdrawAction(ctx).Withdraw(params)
	if internalErr != nil {
		converter.ConvertError(w, internalErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}
