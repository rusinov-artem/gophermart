package error

type ValidationError struct {
	Fields map[string][]string
}

func (e *ValidationError) Error() string {
	return "params has invalid fields"
}
