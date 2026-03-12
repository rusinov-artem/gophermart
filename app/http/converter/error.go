package converter

import (
	"net/http"

	appError "github.com/rusinov-artem/gophermart/app/error"
)

func ConvertError2(w http.ResponseWriter, err *appError.Error) {
	switch e := err.Error.(type) {
	case *appError.InternalError:
		NewInternalErrorConverter(e).Flush(w)
		return
	case *appError.ValidationError:
		NewValidationErrorConverter(e).Flush(w)
		return
	}
}

func ConvertError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *appError.InternalError:
		NewInternalErrorConverter(e).Flush(w)
		return
	case *appError.ValidationError:
		NewValidationErrorConverter(e).Flush(w)
		return
	}
}
