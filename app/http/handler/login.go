package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/app/http/converter"
)

type LoginAction interface {
	Validate(params dto.LoginParams) *appError.ValidationError
	Login(params dto.LoginParams) (string, *appError.InternalError)
}

type loginEnvelop struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (e loginEnvelop) Params() dto.LoginParams {
	return dto.LoginParams{
		Login:    e.Login,
		Password: e.Password,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.Context(r.Context())
	defer cancelFN()

	fmt.Println(ctx)

	body, _ := io.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")

	var envelop loginEnvelop
	err := json.Unmarshal(body, &envelop)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"msg":"unable to unmarshal body"}`))

		return
	}

	action := h.LoginAction(ctx)
	validationErr := action.Validate(envelop.Params())
	if validationErr != nil {
		converter.ConvertError(w, validationErr)
		return
	}

	token, internalError := action.Login(envelop.Params())
	if internalError != nil {
		converter.ConvertError(w, internalError)
		return
	}

	w.Header().Set("Authorization", token)

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}
