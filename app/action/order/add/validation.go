package add

import (
	"fmt"
	"strconv"

	"github.com/rusinov-artem/gophermart/app/action/order"
)

func (t *Action) Validate(orderNr string) error {
	if orderNr == "" {
		return fmt.Errorf("empty order")
	}

	if !t.rule.Match([]byte(orderNr)) {
		return fmt.Errorf("orderNr has invalid format")
	}

	v := order.LuhnCheckDigit(orderNr[:len(orderNr)-1])
	l := orderNr[len(orderNr)-1]
	lv, _ := strconv.Atoi(string(l))

	if v != lv {
		return fmt.Errorf("invalid checksum")
	}

	return nil
}
