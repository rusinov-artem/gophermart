package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type RegisterTestSuite struct {
	suite.Suite
	handler *Handler
	action  *MockRegisterAction
}

func Test_Register(t *testing.T) {
	suite.Run(t, &RegisterTestSuite{})
}

func (s *RegisterTestSuite) SetupTest() {
	s.handler = New()
	s.action = &MockRegisterAction{}
	s.handler.RegisterAction = func(ctx context.Context) RegisterAction {
		return s.action
	}
}

func (s *RegisterTestSuite) Test_EmptyRequest() {
	req := s.req(nil)

	resp := s.handle(req)

	s.Equal(http.StatusBadRequest, resp.Code)
}

func (s *RegisterTestSuite) Test_InvalidParams() {
	req := s.req(bytes.NewBufferString(`{
		"login": "login",
		"password": "password"
	}`))

	s.action.validationError = &appError.ValidationError{Fields: map[string][]string{
		"login": {"invalid login"},
	}}

	resp := s.handle(req)

	s.Equal(http.StatusBadRequest, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
	s.Equal("password", s.action.validationParams.Password)
	s.Equal("login", s.action.validationParams.Login)

	s.JSONEq(
		`{"error":"validation","msg":"some fields not valid","fields":{"login":["invalid login"]}}`,
		resp.Body.String(),
	)
}

func (s *RegisterTestSuite) Test_RegistrationInternalError() {
	req := s.req(bytes.NewBufferString(`{
		"login": "login",
		"password": "password"
	}`))

	s.action.registrationError = &appError.InternalError{
		InnerError: fmt.Errorf("database error"),
		Msg:        "temporary unable to register user",
	}

	resp := s.handle(req)
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
	s.Equal("password", s.action.registrationParams.Password)
	s.Equal("login", s.action.registrationParams.Login)

	s.JSONEq(
		`{"error":"internal","msg":"temporary unable to register user"}`,
		resp.Body.String(),
	)
}

func (s *RegisterTestSuite) Test_LoginAlreadyInUse() {
	req := s.req(bytes.NewBufferString(`{
		"login": "login",
		"password": "password"
	}`))

	s.action.registrationError = &appError.InternalError{
		Msg:  "login already in use",
		Code: appError.LoginAlreadyInUse,
	}

	resp := s.handle(req)
	s.Equal(http.StatusConflict, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))

	s.JSONEq(
		`{"error":"internal","msg":"login already in use"}`,
		resp.Body.String(),
	)
}

func (s *RegisterTestSuite) Test_RegistrationSuccess() {
	req := s.req(bytes.NewBufferString(`{
		"login": "login",
		"password": "password"
	}`))

	s.action.token = "auth token"

	resp := s.handle(req)
	s.Equal(http.StatusOK, resp.Code)
	s.Equal("application/json", resp.Header().Get("Content-Type"))
	s.Equal(s.action.token, resp.Header().Get("Authorization"))
}

func (s *RegisterTestSuite) handle(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.Register(resp, req)
	return resp
}

func (s *RegisterTestSuite) req(body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/api/user/register", body)
}

type MockRegisterAction struct {
	validationError  *appError.ValidationError
	validationParams *dto.RegisterParams

	registrationError  *appError.InternalError
	registrationParams *dto.RegisterParams

	token string
}

func (m *MockRegisterAction) Validate(params dto.RegisterParams) *appError.ValidationError {
	m.validationParams = &params
	return m.validationError
}

func (m *MockRegisterAction) Register(params dto.RegisterParams) (string, *appError.InternalError) {
	m.registrationParams = &params
	return m.token, m.registrationError
}
