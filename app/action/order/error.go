package order

import "fmt"

type NotFoundErr struct {
	OrderNr string
}

func (e *NotFoundErr) Error() string {
	return fmt.Sprintf("order not found: %s", e.OrderNr)
}
