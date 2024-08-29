package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Number(t *testing.T) {
	assert.Equal(t, 0, LuhnCheckDigit("0"))
}
