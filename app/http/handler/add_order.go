package handler

import (
	"io"
	"net/http"
	"strings"

	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/app/http/converter"
)

type AddOrderAction interface {
	Validate(orderNr string) error
	SaveOrder(login, orderNr string) *appError.InternalError
}

type AuthService interface {
	Auth(token string) (string, error)
}

func (h *Handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.Context(r.Context())
	defer cancelFN()

	body, _ := io.ReadAll(r.Body)
	orderNr := strings.TrimSpace(string(body))
	if orderNr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := h.AuthService(ctx)
	login, err := auth.Auth(getToken(r))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	action := h.AddOrderAction(ctx)

	err = action.Validate(orderNr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	internalErr := action.SaveOrder(login, orderNr)
	if internalErr != nil {
		converter.ConvertError(w, internalErr)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func getToken(r *http.Request) string {
	if v := r.Header.Get("Authorization"); v != "" {
		return v
	}

	cookie, err := r.Cookie("Authorization")
	if err != nil {
		return ""
	}

	return cookie.Value
}
