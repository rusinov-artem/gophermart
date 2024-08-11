package error

import "testing"

func Test_error(t *testing.T) {
	i := InternalError{}
	_ = i.Error()

	v := ValidationError{}
	_ = v.Error()
}
