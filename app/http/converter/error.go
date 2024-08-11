package converter

import (
	"net/http"

	appError "github.com/rusinov-artem/gophermart/app/error"
)

func ConvertError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *appError.InternalError:
		if e.Code == appError.LoginAlreadyInUse {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write(NewInternalErrorConverter(e).Byte())

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(NewInternalErrorConverter(e).Byte())

		return
	case *appError.ValidationError:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(NewValidationErrorConverter(e).Byte())

		return
	}
}
