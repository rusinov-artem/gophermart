package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/app/http/converter"
)

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
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(converter.NewValidationErrorConverter(validationErr).Byte())
		return
	}

	token, registrationErr := action.Register(envelop.Params())
	if registrationErr != nil {
		if registrationErr.Code == appError.LoginAlreadyInUse {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write(converter.NewInternalErrorConverter(registrationErr).Byte())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(converter.NewInternalErrorConverter(registrationErr).Byte())
		return
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)

}
