package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_err(t *testing.T) {
	err := NotFoundErr{}
	assert.NotEmpty(t, err.Error())
}
