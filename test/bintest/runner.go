package bintest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func Run(t *testing.T, s suite.TestingSuite) {
	if _, err := os.Stat("./app"); err != nil {
		t.Skip("no file to test")
	}
	suite.Run(t, s)

}
