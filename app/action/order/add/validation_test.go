package add

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidationTestSuite struct {
	suite.Suite
	action *Action
}

func Test_Validation(t *testing.T) {
	suite.Run(t, &ValidationTestSuite{})
}

func (s *ValidationTestSuite) SetupTest() {
	s.action = New(nil, nil)
}

func (s *ValidationTestSuite) Test_EmptyOrderIsInvalid() {
	s.assertInvalid("")
	s.assertInvalid("L")
	s.assertInvalid("1 2")
	s.assertValid("142")
}

func (s *ValidationTestSuite) assertValid(orderNr string) {
	s.NoError(s.action.Validate(orderNr))
}

func (s *ValidationTestSuite) assertInvalid(orderNr string) {
	s.Error(s.action.Validate(orderNr))
}
