package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
	"github.com/rusinov-artem/gophermart/test/utils/logger"
)

type LoginTestSuite struct {
	suite.Suite
	logger  *zap.Logger
	logs    *logger.Logs
	ctx     context.Context
	action  *loginAction
	handler *Handler
}

func Test_LoginHandler(t *testing.T) {
	suite.Run(t, &LoginTestSuite{})
}

func (s *LoginTestSuite) SetupTest() {
	var err error
	s.ctx = context.Background()
	s.Require().NoError(err)

	var logs *bytes.Buffer
	s.logger, logs = logger.SpyLogger()
	s.logs = logger.NewLogs(s.T(), logs)

	s.action = &loginAction{}
	s.handler = New()
	s.handler.LoginAction = func(ctx context.Context) LoginAction {
		return s.action
	}
}

func (s *LoginTestSuite) Test_ErrorIfEmptyBody() {
	req := s.req("")

	resp := s.do(req)

	s.Equal(http.StatusBadRequest, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
	s.Equal(`{"msg":"unable to unmarshal body"}`, resp.Body.String())
}

func (s *LoginTestSuite) Test_ErrorInvalidParams() {
	req := s.req("{}")

	s.action.validationErr = &appError.ValidationError{
		Fields: map[string][]string{"login": {"login is required"}},
	}
	resp := s.do(req)

	s.Equal(http.StatusBadRequest, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
}

func (s *LoginTestSuite) Test_InvalidLoginPassword() {
	req := s.req(`{"login": "test", "password": "test"}`)

	s.action.loginErr = &appError.InternalError{
		InnerError: nil,
		Msg:        "invalid login/password",
		Code:       appError.InvalidCredentials,
	}
	resp := s.do(req)

	s.Equal(http.StatusUnauthorized, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
}

func (s *LoginTestSuite) Test_InternalError() {
	req := s.req(`{"login": "test", "password": "test"}`)

	s.action.loginErr = &appError.InternalError{
		InnerError: nil,
		Msg:        "service temporary unavailable",
		Code:       appError.ServiceUnavailable,
	}
	resp := s.do(req)

	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
}

func (s *LoginTestSuite) Test_SuccessLogin() {
	req := s.req(`{"login": "test", "password": "test"}`)

	s.action.token = "user token"

	resp := s.do(req)

	s.Equal(http.StatusOK, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
	s.Equal(`user token`, resp.Header().Get("Authorization"))
}

func (s *LoginTestSuite) req(body string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBufferString(body))
}

func (s *LoginTestSuite) do(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.Login(resp, req)
	return resp
}

type loginAction struct {
	validateParams *dto.LoginParams
	validationErr  *appError.ValidationError

	loginParams *dto.LoginParams
	loginErr    *appError.InternalError
	token       string
}

func (l *loginAction) Login(params dto.LoginParams) (string, *appError.InternalError) {
	l.loginParams = &params
	return l.token, l.loginErr
}

func (l *loginAction) Validate(params dto.LoginParams) *appError.ValidationError {
	l.validateParams = &params
	return l.validationErr
}
