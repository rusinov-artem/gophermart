package converter

import (
	"net/http"

	appError "github.com/rusinov-artem/gophermart/app/error"
)

func ConvertError(w http.ResponseWriter, err any) {
	switch e := err.(type) {
	case *appError.InternalError:
		NewInternalErrorConverter(e).Flush(w)
		return
	case *appError.ValidationError:
		NewValidationErrorConverter(e).Flush(w)
		return
	}
}
