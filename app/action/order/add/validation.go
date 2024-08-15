package add

import (
	"github.com/rusinov-artem/gophermart/app/action/order"
)

func (t *Action) Validate(orderNr string) error {
	return order.ValidateOrderNr(orderNr)
}
