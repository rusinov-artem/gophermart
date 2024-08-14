package get

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GetBalanceActionTestSuite struct {
	suite.Suite
}

func Test_GetBalanceAction(t *testing.T) {
	suite.Run(t, &GetBalanceActionTestSuite{})
}

func (s *GetBalanceActionTestSuite) SetupTest() {

}

func (s *GetBalanceActionTestSuite) Test_() {}
