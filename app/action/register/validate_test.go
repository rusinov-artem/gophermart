package register

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type ValidationSuite struct {
	suite.Suite
	action  *Register
	storage *storage
	logger  *zap.Logger
	logs    *logger.Logs
}

func Test_Validation(t *testing.T) {
	suite.Run(t, &ValidationSuite{})
}

func (s *ValidationSuite) SetupTest() {
	var logs *bytes.Buffer
	s.logger, logs = logger.SpyLogger()
	s.logs = logger.NewLogs(s.T(), logs)
	s.storage = &storage{}
	s.action = New(s.storage, s.logger, nil)
}

func (s *ValidationSuite) Test_ErrorIfThereIsNoParams() {
	err := s.action.Validate(dto.RegisterParams{})

	s.Error(err)
	s.Equal(err.Fields["password"][0], "password is required")
	s.Equal(err.Fields["loginToSave"][0], "loginToSave is required")
}

func (s *ValidationSuite) Test_ErrorIfPasswordToShort() {
	err := s.action.Validate(dto.RegisterParams{
		Login:    "loginToSave",
		Password: "123",
	})

	s.Require().Error(err)
	s.Equal(err.Fields["password"][0], "password is too short")
}

func (s *ValidationSuite) Test_SuccessValidation() {
	err := s.action.Validate(dto.RegisterParams{
		Login:    "login",
		Password: "nice password to use",
	})
	s.Nil(err)
}

type storage struct {
	isLoginExists      bool
	isLoginExistsError error
	loginToCheck       string

	loginToSave string
	password    string
	saveError   error

	token         string
	loginForToken string
	addTokenErr   error
}

func (s *storage) SaveUser(login, password string) error {
	s.loginToSave = login
	s.password = password
	return s.saveError
}

func (s *storage) IsLoginExists(login string) (bool, error) {
	s.loginToCheck = login
	return s.isLoginExists, s.isLoginExistsError
}

func (s *storage) AddToken(login, token string) error {
	s.loginForToken = login
	s.token = token
	return s.addTokenErr
}
