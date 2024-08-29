package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type GetBalanceHandlerTestSuite struct {
	suite.Suite
	handler *Handler
	auth    *authService
	action  *getBalanceAction
}

func Test_BalanceHandler(t *testing.T) {
	suite.Run(t, &GetBalanceHandlerTestSuite{})
}

func (s *GetBalanceHandlerTestSuite) SetupTest() {
	s.action = &getBalanceAction{}
	s.auth = &authService{}
	s.handler = New()
	s.handler.AuthService = func(ctx context.Context) AuthService {
		return s.auth
	}

	s.handler.GetBalanceAction = func(ctx context.Context) GetBalanceAction {
		return s.action
	}
}

func (s *GetBalanceHandlerTestSuite) Test_Unauthorized() {
	s.auth.autErr = fmt.Errorf("user not authorized")

	resp := s.do(s.req())

	s.Equal(http.StatusUnauthorized, resp.Code)
}

func (s *GetBalanceHandlerTestSuite) Test_UnableToGetBalance() {
	login := "login"

	s.auth.login = login

	s.action.err = &appError.InternalError{
		Msg:  "service temporary unavailable",
		Code: appError.ServiceUnavailable,
	}

	resp := s.do(s.req())

	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Equal(login, s.action.login)
}

func (s *GetBalanceHandlerTestSuite) Test_CanGetBalance() {
	login := "login"

	s.auth.login = login

	s.action.balance = dto.Balance{
		Current:   1000,
		Withdrawn: 500,
	}

	resp := s.do(s.req())

	s.Equal(http.StatusOK, resp.Code)
	s.Equal(login, s.action.login)
	fmt.Println(resp.Body.String())
	s.JSONEq(`
	  {
		  "current": 1000,
		  "withdrawn": 500
	  }`,
		resp.Body.String(),
	)
}

func (s *GetBalanceHandlerTestSuite) req() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
}

func (s *GetBalanceHandlerTestSuite) do(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.GetBalance(resp, req)
	return resp
}

type getBalanceAction struct {
	login   string
	err     *appError.InternalError
	balance dto.Balance
}

func (t *getBalanceAction) GetBalance(login string) (dto.Balance, *appError.InternalError) {
	t.login = login
	return t.balance, t.err
}
