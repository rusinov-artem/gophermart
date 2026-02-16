package converter

import (
	"encoding/json"
	"net/http"

	appError "github.com/rusinov-artem/gophermart/app/error"
)

type jsonValidationError struct {
	Error  string              `json:"error"`
	Msg    string              `json:"msg"`
	Fields map[string][]string `json:"fields"`
}

type ValidationErrorConverter struct {
	err *appError.ValidationError
}

func NewValidationErrorConverter(err *appError.ValidationError) *ValidationErrorConverter {
	return &ValidationErrorConverter{
		err: err,
	}
}

func (t *ValidationErrorConverter) Byte() []byte {
	e := jsonValidationError{
		Error:  "validation",
		Msg:    "some fields not valid",
		Fields: t.err.Fields,
	}

	b, _ := json.Marshal(e)

	return b
}

func (t *ValidationErrorConverter) Flush(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(t.Byte())
}
