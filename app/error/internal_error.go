package error

type InternalError struct {
	InnerError error
	Msg        string
	Code       Code
}

func (t *InternalError) Error() string {
	return t.Msg
}
