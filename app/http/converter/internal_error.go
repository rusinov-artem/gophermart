package converter

import (
	"encoding/json"
	"net/http"

	appError "github.com/rusinov-artem/gophermart/app/error"
)

var statusMap = map[appError.Code]int{
	appError.LoginAlreadyInUse:  http.StatusConflict,
	appError.InvalidCredentials: http.StatusUnauthorized,
	appError.OrderNrExists:      http.StatusOK,
	appError.BadOrderOwnership:  http.StatusConflict,
	appError.NoOrdersFound:      http.StatusNoContent,
	appError.NotEnoughPoints:    http.StatusPaymentRequired,
	appError.NoWithdrawals:      http.StatusNoContent,
}

type InternalErrorConverter struct {
	err *appError.InternalError
}

type jsonInternalError struct {
	Error string `json:"error"`
	Msg   string `json:"msg"`
}

func NewInternalErrorConverter(err *appError.InternalError) *InternalErrorConverter {
	return &InternalErrorConverter{err: err}
}

func (t *InternalErrorConverter) Byte() []byte {
	e := jsonInternalError{
		Error: "internal",
		Msg:   t.err.Error(),
	}

	b, _ := json.Marshal(e)

	return b
}

func (t *InternalErrorConverter) Flush(w http.ResponseWriter) {
	if status, ok := statusMap[t.err.Code]; ok {
		w.WriteHeader(status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, _ = w.Write(t.Byte())
}
