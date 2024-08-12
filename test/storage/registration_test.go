package storage

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/storage"
	"github.com/rusinov-artem/gophermart/test"
)

type RegistrationStorageTestSuite struct {
	suite.Suite
	pool    *pgxpool.Pool
	ctx     context.Context
	storage *storage.RegistrationStorage
}

func Test_RegistrationStorage(t *testing.T) {
	suite.Run(t, &RegistrationStorageTestSuite{})
}

func (s *RegistrationStorageTestSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()
	dsn := test.CreateTestDB("test_registration_storage")
	s.pool, err = pgxpool.New(context.Background(), dsn)
	s.Require().NoError(err)
}

func (s *RegistrationStorageTestSuite) SetupTest() {
	s.storage = storage.NewRegistrationStorage(s.ctx, s.pool)
}

func (s *RegistrationStorageTestSuite) Test_CanSaveUser() {
	login := "TestLogin001"
	password := "password"
	isExists, err := s.storage.IsLoginExists(login)
	s.Require().NoError(err)
	s.False(isExists)

	err = s.storage.SaveUser(login, password)
	s.Require().NoError(err)

	isExists, err = s.storage.IsLoginExists(login)
	s.Require().NoError(err)
	s.True(isExists)
}

func (s *RegistrationStorageTestSuite) Test_CanAddToken() {
	login := "TestLogin002"
	password := "password"
	err := s.storage.SaveUser(login, password)
	s.Require().NoError(err)

	token := "some token"
	err = s.storage.AddToken(login, token)
	s.Require().NoError(err)

	s.AssertToken(token, login)
}

func (s *RegistrationStorageTestSuite) AssertToken(token, login string) {
	s.T().Helper()

	sqlStr := `SELECT token, login FROM auth_token WHERE token = $1`
	rows, err := s.pool.Query(s.ctx, sqlStr, token)
	s.Require().NoError(err)
	s.Require().True(rows.Next())

	var found struct {
		token string
		login string
	}

	s.Require().NoError(rows.Scan(&found.token, &found.login))

	s.Equal(token, found.token)
	s.Equal(login, found.login)
}
