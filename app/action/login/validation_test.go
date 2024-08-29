package login

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/dto"
)

type ValidationTestSuite struct {
	suite.Suite
	action *Login
}

func Test_Validation(t *testing.T) {
	suite.Run(t, &ValidationTestSuite{})
}

func (s *ValidationTestSuite) SetupTest() {
	s.action = New(nil, nil, nil)
}

func (s *ValidationTestSuite) Test_ErrorIfParamsEmpty() {
	err := s.action.Validate(dto.LoginParams{})
	s.NotNil(err)
	s.Equal("login is required", err.Fields["login"][0])
	s.Equal("password is required", err.Fields["password"][0])
}

func (s *ValidationTestSuite) Test_SuccessValidation() {
	err := s.action.Validate(dto.LoginParams{
		Login:    "login",
		Password: "password",
	})
	s.Nil(err)
}
