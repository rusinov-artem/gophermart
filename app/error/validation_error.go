package error

type ValidationError struct {
	Fields map[string][]string
}

func (t *ValidationError) appError() {}

func (e *ValidationError) Error() string {
	return "params has invalid fields"
}
