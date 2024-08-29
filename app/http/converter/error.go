package converter

import (
	"net/http"

	appError "github.com/rusinov-artem/gophermart/app/error"
)

func ConvertError(w http.ResponseWriter, err interface{}) {
	switch e := err.(type) {
	case *appError.InternalError:
		NewInternalErrorConverter(e).Flush(w)
		return
	case *appError.ValidationError:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(NewValidationErrorConverter(e).Byte())
		return
	}
}
