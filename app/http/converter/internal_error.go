package converter

import (
	"encoding/json"

	appError "github.com/rusinov-artem/gophermart/app/error"
)

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
