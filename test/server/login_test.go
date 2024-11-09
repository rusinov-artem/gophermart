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

	"github.com/rusinov-artem/gophermart/app/http/handler"
	"github.com/rusinov-artem/gophermart/test"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type LoginTestSuite struct {
	suite.Suite
	logger  *zap.Logger
	logs    *logger.Logs
	handler *handler.Handler
	ctx     context.Context
	dsn     string
	pool    *pgxpool.Pool
}

func Test_Login(t *testing.T) {
	suite.Run(t, &LoginTestSuite{})
}

func (s *LoginTestSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()
	s.dsn = test.CreateTestDB("test_login")
	s.pool, err = pgxpool.New(context.Background(), s.dsn)
	s.Require().NoError(err)
}

func (s *LoginTestSuite) SetupTest() {
	var logs *bytes.Buffer
	s.logger, logs = logger.SpyLogger()
	s.logs = logger.NewLogs(s.T(), logs)

	s.handler = handler.New()
}

func (s *LoginTestSuite) req(body string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(body))
}

func (s *LoginTestSuite) do(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.Login(resp, req)
	return resp
}
