package order

type InvalidOrderNrFormat struct{}

func (t *InvalidOrderNrFormat) Error() string {
	return "order number is invalid"
}
