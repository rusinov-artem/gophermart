package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rusinov-artem/gophermart/app/action/order"
	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type WithdrawHandlerTestSuite struct {
	suite.Suite
	handler *Handler
	auth    *authService
	action  *withdrawAction
}

func Test_WithdrawHandler(t *testing.T) {
	suite.Run(t, &WithdrawHandlerTestSuite{})
}

func (s *WithdrawHandlerTestSuite) SetupTest() {
	s.auth = &authService{}
	s.action = &withdrawAction{}
	s.handler = New()
	s.handler.AuthService = func(ctx context.Context) AuthService {
		return s.auth
	}

	s.handler.WithdrawAction = func(ctx context.Context) WithdrawAction {
		return s.action
	}
}

func (s *WithdrawHandlerTestSuite) Test_UnableToUnmarshalBody() {
	req := s.req("}InvalidJson{")
	resp := s.do(req)

	s.Equal(http.StatusBadRequest, resp.Code)
}

func (s *WithdrawHandlerTestSuite) Test_InvalidOrderNr() {
	req := s.req(`{"number":"invalid", "sum":42}`)
	resp := s.do(req)

	s.Equal(http.StatusUnprocessableEntity, resp.Code)
}

func (s *WithdrawHandlerTestSuite) Test_UnauthorizedUser() {
	s.auth.autErr = fmt.Errorf("unknown token")

	orderNr := order.OrderNr()
	resp := s.do(s.validReq(orderNr, 42))

	s.Equal(http.StatusUnauthorized, resp.Code)
}

func (s *WithdrawHandlerTestSuite) Test_UnableToWithdraw() {
	login := "login"
	s.auth.login = login

	s.action.err = &appError.InternalError{
		InnerError: fmt.Errorf("db error"),
		Msg:        "service temporary unavailable",
		Code:       appError.ServiceUnavailable,
	}

	orderNr := order.OrderNr()
	sum := float32(42.42)
	resp := s.do(s.validReq(orderNr, sum))

	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Equal(login, s.action.params.Login)
	s.Equal(orderNr, s.action.params.OrderNr)
	s.InDelta(sum, s.action.params.Sum, 0.0001)
}

func (s *WithdrawHandlerTestSuite) Test_NotEnoughPoints() {
	login := "login"
	s.auth.login = login

	s.action.err = &appError.InternalError{
		Msg:  "not enough points",
		Code: appError.NotEnoughPoints,
	}

	orderNr := order.OrderNr()
	sum := float32(42.42)
	resp := s.do(s.validReq(orderNr, sum))

	s.Equal(http.StatusPaymentRequired, resp.Code)
	s.Equal(login, s.action.params.Login)
	s.Equal(orderNr, s.action.params.OrderNr)
	s.InDelta(sum, s.action.params.Sum, 0.0001)
}

func (s *WithdrawHandlerTestSuite) Test_Success() {
	login := "login"
	s.auth.login = login

	orderNr := order.OrderNr()
	sum := float32(42.42)
	resp := s.do(s.validReq(orderNr, sum))

	s.Equal(http.StatusOK, resp.Code)
	s.Equal(login, s.action.params.Login)
	s.Equal(orderNr, s.action.params.OrderNr)
	s.InDelta(sum, s.action.params.Sum, 0.0001)
}

func (s *WithdrawHandlerTestSuite) validReq(orderNr string, sum float32) *http.Request {
	return s.req(fmt.Sprintf(`{"number":"%s", "sum":%f}`, orderNr, sum))
}

func (s *WithdrawHandlerTestSuite) req(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBufferString(body))
	req.Header.Add("Content-Type", "application/json")
	return req
}

func (s *WithdrawHandlerTestSuite) do(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.Withdraw(resp, req)
	return resp
}

type withdrawAction struct {
	params dto.WithdrawParams
	err    *appError.InternalError
}

func (t *withdrawAction) Withdraw(params dto.WithdrawParams) *appError.InternalError {
	t.params = params
	return t.err
}
