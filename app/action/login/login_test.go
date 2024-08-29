package login

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

type LoginTestSuite struct {
	suite.Suite
	action    *Login
	storage   *storage
	logs      *bytes.Buffer
	logger    *zap.Logger
	generator *generator
}

func Test_Login(t *testing.T) {
	suite.Run(t, &LoginTestSuite{})
}

func (s *LoginTestSuite) SetupTest() {
	s.storage = &storage{}
	s.generator = &generator{}
	s.logger, s.logs = logger.SpyLogger()
	s.action = New(s.storage, s.logger, s.generator)
}

func (s *LoginTestSuite) Test_UnableToFundUser() {
	params := dto.LoginParams{
		Login:    "login",
		Password: "password",
	}

	s.storage.findUserErr = fmt.Errorf("unable to find user")

	_, err := s.action.Login(params)
	s.NotNil(err)
	s.Equal("login", s.storage.loginToSearch)
}

func (s *LoginTestSuite) Test_UserNotFound() {
	params := dto.LoginParams{
		Login:    "login",
		Password: "password",
	}

	s.storage.findUserErr = &UserNotFoundErr{Login: "login"}

	_, err := s.action.Login(params)
	s.NotNil(err)
	s.Equal("login", s.storage.loginToSearch)
	s.Equal(appError.InvalidCredentials, err.Code)
}

func (s *LoginTestSuite) Test_InvalidPassword() {
	params := dto.LoginParams{
		Login:    "login",
		Password: "password",
	}

	s.storage.foundUser = dto.User{
		Login:        "user",
		PasswordHash: "invalid hash",
	}

	s.action.CheckPasswordHash = func(_, _ string) bool {
		return false
	}

	_, err := s.action.Login(params)
	s.NotNil(err)
	s.Equal("login", s.storage.loginToSearch)
	s.Equal(appError.InvalidCredentials, err.Code)
}

func (s *LoginTestSuite) Test_UnableToSaveToken() {
	params := dto.LoginParams{
		Login:    "login",
		Password: "password",
	}

	s.storage.foundUser = dto.User{
		Login:        "login",
		PasswordHash: "valid hash",
	}

	s.action.CheckPasswordHash = func(_, _ string) bool {
		return true
	}

	s.generator.token = "user token"

	s.storage.addTokenErr = fmt.Errorf("unable to save token")

	_, err := s.action.Login(params)
	s.NotNil(err)
	s.Equal("login", s.storage.loginToSearch)
	s.Equal("login", s.storage.loginForToken)
	s.Equal(s.generator.token, s.storage.tokenToSave)
	s.Equal(appError.ServiceUnavailable, err.Code)
}

func (s *LoginTestSuite) Test_LoginSuccess() {
	params := dto.LoginParams{
		Login:    "login",
		Password: "password",
	}

	s.storage.foundUser = dto.User{
		Login:        "login",
		PasswordHash: "valid hash",
	}

	s.action.CheckPasswordHash = func(_, _ string) bool {
		return true
	}

	s.generator.token = "user token"

	token, err := s.action.Login(params)
	s.Nil(err)
	s.Equal("login", s.storage.loginToSearch)
	s.Equal("login", s.storage.loginForToken)
	s.Equal(s.generator.token, s.storage.tokenToSave)
	s.Equal(s.generator.token, token)
}

func (s *LoginTestSuite) Test_Err() {
	err := &UserNotFoundErr{}
	s.NotEmpty(err.Error())
}

type storage struct {
	findUserErr   error
	loginToSearch string
	foundUser     dto.User

	loginForToken string
	tokenToSave   string
	addTokenErr   error
}

func (s *storage) FindUser(login string) (dto.User, error) {
	s.loginToSearch = login
	return s.foundUser, s.findUserErr
}

func (s *storage) AddToken(login, token string) error {
	s.loginForToken = login
	s.tokenToSave = token
	return s.addTokenErr
}

type generator struct {
	token string
}

func (g *generator) Generate() string {
	return g.token
}
