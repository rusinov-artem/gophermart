package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type GetWithdrawalsHandlerTestSuite struct {
	suite.Suite
	handler *Handler
	auth    *authService
	action  *getWithdrawalsAction
}

func Test_GetWithdrawalsHandler(t *testing.T) {
	suite.Run(t, &GetWithdrawalsHandlerTestSuite{})
}

func (s *GetWithdrawalsHandlerTestSuite) SetupTest() {
	s.auth = &authService{}
	s.handler = New()
	s.handler.AuthService = func(ctx context.Context) AuthService {
		return s.auth
	}

	s.action = &getWithdrawalsAction{}
	s.handler.GetWithdrawalsAction = func(ctx context.Context) GetWithdrawalsAction {
		return s.action
	}
}

func (s *GetWithdrawalsHandlerTestSuite) Test_UserUnauthorized() {
	s.auth.autErr = fmt.Errorf("unknown user")
	resp := s.do(s.req())
	s.Equal(http.StatusUnauthorized, resp.Code)
}

func (s *GetWithdrawalsHandlerTestSuite) Test_UnableToGetWithdrawals() {
	s.auth.login = "login"
	s.action.err = &appError.InternalError{
		Code: appError.ServiceUnavailable,
	}

	resp := s.do(s.req())
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Equal(s.auth.login, s.action.login)
}

func (s *GetWithdrawalsHandlerTestSuite) Test_NoWithdrawals() {
	s.auth.login = "login"
	s.action.err = &appError.InternalError{
		Code: appError.NoWithdrawals,
	}

	resp := s.do(s.req())
	s.Equal(http.StatusNoContent, resp.Code)
	s.Equal(s.auth.login, s.action.login)
}

func (s *GetWithdrawalsHandlerTestSuite) Test_CanGetWithdrawals() {
	s.auth.login = "login"
	dt, _ := time.Parse(time.RFC3339, "2020-12-09T16:09:57+03:00")
	s.action.withdrawals = []dto.Withdrawal{
		{
			OrderNr:     "2377225624",
			Sum:         500,
			ProcessedAt: dt,
		},
	}

	resp := s.do(s.req())
	s.Equal(http.StatusOK, resp.Code)
	s.Equal(s.auth.login, s.action.login)
	s.JSONEq(`
		  [
			  {
				  "order": "2377225624",
				  "sum": 500,
				  "processed_at": "2020-12-09T16:09:57+03:00"
			  }
		  ]`, resp.Body.String())
}

func (s *GetWithdrawalsHandlerTestSuite) req() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
}

func (s *GetWithdrawalsHandlerTestSuite) do(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.GetWithdrawals(resp, req)
	return resp
}

type getWithdrawalsAction struct {
	login       string
	withdrawals []dto.Withdrawal
	err         *appError.InternalError
}

func (s *getWithdrawalsAction) GetWithdrawals(login string) ([]dto.Withdrawal, *appError.InternalError) {
	s.login = login
	return s.withdrawals, s.err
}
