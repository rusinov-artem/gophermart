package register

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type RegisterTestSuite struct {
	suite.Suite
	storage        *storage
	action         *Register
	logger         *zap.Logger
	logs           *logger.Logs
	tokenGenerator *tokenGenerator
}

func Test_Register(t *testing.T) {
	suite.Run(t, &RegisterTestSuite{})
}

func (s *RegisterTestSuite) SetupTest() {
	var logs *bytes.Buffer
	s.logger, logs = logger.SpyLogger()
	s.logs = logger.NewLogs(s.T(), logs)
	s.storage = &storage{}
	s.tokenGenerator = &tokenGenerator{}
	s.action = New(s.storage, s.logger, s.tokenGenerator)
}

func (s *RegisterTestSuite) Test_UnableToCheckLoginExists() {
	s.storage.isLoginExistsError = fmt.Errorf("database error")

	token, err := s.action.Register(dto.RegisterParams{
		Login:    "loginToSave",
		Password: "password",
	})

	s.Error(err)
	s.Empty(token)

	s.logs.Contains("error", "unable to check loginToSave exists", "database error")
}

func (s *RegisterTestSuite) Test_ErrorIfUserAlreadyExists() {
	s.storage.isLoginExists = true

	token, err := s.action.Register(dto.RegisterParams{
		Login:    "loginToSave",
		Password: "password",
	})

	s.Error(err)
	s.Equal("loginToSave", s.storage.loginToCheck)
	s.Equal(err.Msg, "login already in use")
	s.Equal(err.Code, appError.LoginAlreadyInUse)
	s.Empty(token)
}

func (s *RegisterTestSuite) Test_ErrorIfUnableToSaveUser() {
	s.storage.saveError = fmt.Errorf("database error")

	token, err := s.action.Register(dto.RegisterParams{
		Login:    "loginToSave",
		Password: "password",
	})

	s.Error(err)
	s.Equal(err.Msg, "unable to save user")
	s.Equal(err.Code, appError.UnableToSaveUser)
	s.Empty(token)

	s.logs.Contains("error", "unable to save user", "database error")
}

func (s *RegisterTestSuite) Test_ErrorIfUnableToAddToken() {
	s.tokenGenerator.token = "some token"
	s.storage.addTokenErr = fmt.Errorf("db error")

	token, err := s.action.Register(dto.RegisterParams{
		Login:    "login",
		Password: "password",
	})

	s.Error(err)
	s.Equal(err.Msg, "unable to save token")
	s.Equal(err.Code, appError.UnableToSaveToken)
	s.Empty(token)

	s.logs.Contains("error", "unable to save token", "db error")
}

func (s *RegisterTestSuite) Test_CanGetToken() {
	s.tokenGenerator.token = "some token"

	token, err := s.action.Register(dto.RegisterParams{
		Login:    "login",
		Password: "password",
	})

	s.Nil(err)
	s.Equal("some token", token)
}

type tokenGenerator struct {
	token string
}

func (t *tokenGenerator) Generate() string {
	return t.token
}
