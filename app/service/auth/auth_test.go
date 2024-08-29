package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	service *Service
	storage *storage
}

func Test_AuthService(t *testing.T) {
	suite.Run(t, &AuthServiceTestSuite{})
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.storage = &storage{}
	s.service = NewService(s.storage)
}

func (s *AuthServiceTestSuite) Test_UnableToFindToken() {
	token := "unknown"
	s.storage.findErr = fmt.Errorf("db error")

	_, err := s.service.Auth(token)
	s.Error(err)
	s.Equal(token, s.storage.tokenToFind)
}

func (s *AuthServiceTestSuite) Test_Success() {
	token := "unknown"
	s.storage.foundLogin = "login"

	login, err := s.service.Auth(token)
	s.NoError(err)
	s.Equal(s.storage.foundLogin, login)
}

type storage struct {
	tokenToFind string
	findErr     error
	foundLogin  string
}

func (s *storage) FindToken(token string) (string, error) {
	s.tokenToFind = token
	return s.foundLogin, s.findErr
}
