package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	appError "github.com/rusinov-artem/gophermart/app/error"
)

type AddOrderHandlerTestSuite struct {
	suite.Suite
	handler     *Handler
	action      *addOrderAction
	authService *authService
}

func Test_AddOrderHandler(t *testing.T) {
	suite.Run(t, &AddOrderHandlerTestSuite{})
}

func (s *AddOrderHandlerTestSuite) SetupTest() {
	s.action = &addOrderAction{}
	s.authService = &authService{}

	s.handler = New()

	s.handler.AuthService = func(ctx context.Context) AuthService {
		return s.authService
	}

	s.handler.AddOrderAction = func(ctx context.Context) AddOrderAction {
		return s.action
	}
}

func (s *AddOrderHandlerTestSuite) Test_EmptyOrderNr() {
	s.authService.autErr = fmt.Errorf("auth error")

	orderNr := ""
	req := s.req(orderNr)
	resp := s.do(req)

	s.Equal(http.StatusBadRequest, resp.Code)
}

func (s *AddOrderHandlerTestSuite) Test_ErrorUnauthorized() {
	s.authService.autErr = fmt.Errorf("auth error")

	orderNr := "orderNr"
	req := s.req(orderNr)
	resp := s.do(req)

	s.Equal(http.StatusUnauthorized, resp.Code)
}

func (s *AddOrderHandlerTestSuite) Test_ErrorIfOrderNumberIsInvalid() {
	s.action.validationErr = fmt.Errorf("invalid orderNr")

	orderNr := "invalid order number"
	req := s.req(orderNr)
	resp := s.do(req)

	s.Equal(http.StatusUnprocessableEntity, resp.Code)
	s.Equal(orderNr, s.action.orderNrToValidate)
}

func (s *AddOrderHandlerTestSuite) Test_OrderAlreadyExists() {
	s.action.saveOrderErr = &appError.InternalError{
		InnerError: nil,
		Msg:        "order already exists",
		Code:       appError.OrderNrExists,
	}

	s.authService.login = "login"

	orderNr := "existingOrder"
	req := s.req(orderNr)
	resp := s.do(req)

	s.Equal(http.StatusOK, resp.Code)
	s.Equal(orderNr, s.action.orderNrToValidate)
	s.Equal(orderNr, s.action.orderNrToSave)
	s.Equal(s.authService.login, s.action.login)
}

func (s *AddOrderHandlerTestSuite) Test_OrderBelongToDifferentUser() {
	s.action.saveOrderErr = &appError.InternalError{
		InnerError: nil,
		Msg:        "order owned by different user",
		Code:       appError.BadOrderOwnership,
	}

	orderNr := "foreignOrderNr"
	req := s.req(orderNr)
	resp := s.do(req)

	s.Equal(http.StatusConflict, resp.Code)
}

func (s *AddOrderHandlerTestSuite) Test_InternalError() {
	s.action.saveOrderErr = &appError.InternalError{
		InnerError: fmt.Errorf("database error"),
		Msg:        "order owned by different user",
		Code:       appError.ServiceUnavailable,
	}

	orderNr := "foreignOrderNr"
	req := s.req(orderNr)
	resp := s.do(req)

	s.Equal(http.StatusInternalServerError, resp.Code)
}

func (s *AddOrderHandlerTestSuite) Test_AddedNewOrder() {
	orderNr := "newOrderNr"
	s.authService.login = "login"

	req := s.req(orderNr)
	resp := s.do(req)

	s.Equal(http.StatusAccepted, resp.Code)
	s.Equal(orderNr, s.action.orderNrToValidate)
	s.Equal(orderNr, s.action.orderNrToSave)
	s.Equal(s.authService.login, s.action.login)
}

func (s *AddOrderHandlerTestSuite) req(number string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(number))
}

func (s *AddOrderHandlerTestSuite) do(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.AddOrder(resp, req)
	return resp
}

type addOrderAction struct {
	validationErr     error
	orderNrToValidate string

	orderNrToSave string
	login         string
	saveOrderErr  *appError.InternalError
}

func (t *addOrderAction) Validate(orderNr string) error {
	t.orderNrToValidate = orderNr
	return t.validationErr
}

func (t *addOrderAction) SaveOrder(login, orderNr string) *appError.InternalError {
	t.orderNrToSave = orderNr
	t.login = login
	return t.saveOrderErr
}

type authService struct {
	tokenToCheck string
	login        string
	autErr       error
}

func (a *authService) Auth(token string) (string, error) {
	a.tokenToCheck = token
	return a.login, a.autErr
}
