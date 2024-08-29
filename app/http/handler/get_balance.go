package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/app/http/converter"
)

type GetBalanceAction interface {
	GetBalance(login string) (dto.Balance, *appError.InternalError)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.Context(r.Context())
	defer cancelFN()

	auth := h.AuthService(ctx)
	login, err := auth.Auth(getToken(r))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	action := h.GetBalanceAction(ctx)
	balance, internalErr := action.GetBalance(login)
	if internalErr != nil {
		converter.ConvertError(w, internalErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(encodeBalance(balance))
}

func encodeBalance(balance dto.Balance) []byte {
	jsonBalance := struct {
		Current   float32 `json:"current"`
		Withdrawn float32 `json:"withdrawn"`
	}{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}

	b, _ := json.Marshal(jsonBalance)
	return b
}
