package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/app/http/converter"
)

const CookiesLifeTime = 3600

type RegisterAction interface {
	Validate(params dto.RegisterParams) *appError.ValidationError
	Register(params dto.RegisterParams) (string, *appError.InternalError)
}

type registerEnvelop struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (e registerEnvelop) Params() dto.RegisterParams {
	return dto.RegisterParams{
		Login:    e.Login,
		Password: e.Password,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancelFN := h.Context(r.Context())
	defer cancelFN()

	w.Header().Set("Content-Type", "application/json")

	data, _ := io.ReadAll(r.Body)
	var envelop registerEnvelop
	err := json.Unmarshal(data, &envelop)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"msg":"unable to unmarshal body"}`))
		return
	}

	action := h.RegisterAction(ctx)
	validationErr := action.Validate(envelop.Params())
	if validationErr != nil {
		converter.ConvertError(w, validationErr)
		return
	}

	token, registrationErr := action.Register(envelop.Params())
	if registrationErr != nil {
		converter.ConvertError(w, registrationErr)
		return
	}

	w.Header().Set("Authorization", token)

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		MaxAge:   CookiesLifeTime,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}
