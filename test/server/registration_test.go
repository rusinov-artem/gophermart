package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/action/register"
	"github.com/rusinov-artem/gophermart/app/crypto"
	"github.com/rusinov-artem/gophermart/app/http/handler"
	appStorage "github.com/rusinov-artem/gophermart/app/storage"
	"github.com/rusinov-artem/gophermart/test"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type RegistrationTestSuite struct {
	suite.Suite
	logger  *zap.Logger
	logs    *logger.Logs
	ctx     context.Context
	dsn     string
	handler *handler.Handler
	pool    *pgxpool.Pool
}

func Test_Registration(t *testing.T) {
	suite.Run(t, &RegistrationTestSuite{})
}

func (s *RegistrationTestSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()
	s.dsn = test.CreateTestDB("test_registration")
	s.pool, err = pgxpool.New(context.Background(), s.dsn)
	s.Require().NoError(err)
}

func (s *RegistrationTestSuite) SetupTest() {
	var logs *bytes.Buffer
	s.logger, logs = logger.SpyLogger()
	s.logs = logger.NewLogs(s.T(), logs)

	s.handler = handler.New()

	s.handler.RegisterAction = func(ctx context.Context) handler.RegisterAction {
		storage := appStorage.NewStorage(ctx, s.pool)
		generator := crypto.NewTokenGenerator()
		return register.New(storage, s.logger, generator)
	}
}

func (s *RegistrationTestSuite) Test_CanRegisterAUser() {
	req := s.registerReq(`{"login":"login", "password": "password"}`)

	resp := s.do(req)

	s.Equal(http.StatusOK, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
}

func (s *RegistrationTestSuite) Test_ErrorIfUserAlreadyExists() {
	req := s.registerReq(`{"login":"alreadyExists", "password": "password"}`)
	_ = s.do(req)

	req = s.registerReq(`{"login":"alreadyExists", "password": "password"}`)
	resp := s.do(req)

	s.Equal(http.StatusConflict, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
}

func (s *RegistrationTestSuite) registerReq(body string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBufferString(body))
}

func (s *RegistrationTestSuite) do(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.Register(resp, req)
	return resp
}
