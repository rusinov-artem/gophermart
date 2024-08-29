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

type ListOrdersHandlerTestSuite struct {
	suite.Suite
	handler *Handler
	auth    *authService
	action  *listOrderAction
}

func Test_ListOrders(t *testing.T) {
	suite.Run(t, &ListOrdersHandlerTestSuite{})
}

func (s *ListOrdersHandlerTestSuite) SetupTest() {
	s.auth = &authService{}
	s.action = &listOrderAction{}
	s.handler = New()
	s.handler.AuthService = func(ctx context.Context) AuthService {
		return s.auth
	}

	s.handler.ListOrdersAction = func(ctx context.Context) ListOrdersAction {
		return s.action
	}
}

func (s *ListOrdersHandlerTestSuite) Test_Unauthorized() {
	s.auth.autErr = fmt.Errorf("token not found")

	resp := s.do(s.req())

	s.Equal(http.StatusUnauthorized, resp.Code)
}

func (s *ListOrdersHandlerTestSuite) Test_UnableToFindOrders() {
	s.action.listOrdersErr = &appError.InternalError{
		InnerError: fmt.Errorf("db error"),
		Msg:        "service unavailable",
		Code:       appError.ServiceUnavailable,
	}

	s.auth.login = "login"

	resp := s.do(s.req())

	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Equal(s.auth.login, s.action.login)
}

func (s *ListOrdersHandlerTestSuite) Test_UserDoesNotHaveOrders() {
	s.action.listOrdersErr = &appError.InternalError{
		InnerError: fmt.Errorf("db error"),
		Msg:        "no orders found",
		Code:       appError.NoOrdersFound,
	}

	s.auth.login = "login"

	resp := s.do(s.req())

	s.Equal(http.StatusNoContent, resp.Code)
	s.Equal(s.auth.login, s.action.login)
}

func (s *ListOrdersHandlerTestSuite) Test_Success() {
	s.action.foundOrders = []dto.OrderListItem{
		{
			OrderNr: "1111111111",
			Status:  "PROCESSED",
			Accrual: 500,
			UploadAt: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2020-12-10T15:15:45+03:00")
				return t
			}(),
		},
		{
			OrderNr: "222222222",
			Status:  "PROCESSING",
			UploadAt: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2020-12-10T15:12:01+03:00")
				return t
			}(),
		},
		{
			OrderNr: "3333333",
			Status:  "INVALID",
			UploadAt: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2020-12-10T15:12:01+03:00")
				return t
			}(),
		},
		{
			OrderNr: "44444444",
			Status:  "NEW",
			UploadAt: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2020-12-10T15:12:01+03:00")
				return t
			}(),
		},
	}

	s.auth.login = "login"

	resp := s.do(s.req())
	s.JSONEq(`
		[
			{
				"number": "1111111111",
				"status": "PROCESSED",
				"accrual": 500,
				"uploaded_at": "2020-12-10T15:15:45+03:00"
			},
			{
				"number": "222222222",
				"status": "PROCESSING",
				"uploaded_at": "2020-12-10T15:12:01+03:00"
			},
			{
				"number": "3333333",
				"status": "INVALID",
				"uploaded_at": "2020-12-10T15:12:01+03:00"
			},
			{
				"number": "44444444",
				"status": "NEW",
				"uploaded_at": "2020-12-10T15:12:01+03:00"
			}
		]
`, resp.Body.String())
	s.Equal(http.StatusOK, resp.Code)
	s.Equal(s.auth.login, s.action.login)
}

func (s *ListOrdersHandlerTestSuite) req() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
}

func (s *ListOrdersHandlerTestSuite) do(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.handler.ListOrders(resp, req)
	return resp
}

type listOrderAction struct {
	listOrdersErr *appError.InternalError
	login         string
	foundOrders   []dto.OrderListItem
}

func (t *listOrderAction) ListOrders(login string) ([]dto.OrderListItem, *appError.InternalError) {
	t.login = login
	return t.foundOrders, t.listOrdersErr
}
