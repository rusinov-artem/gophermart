package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetToken(t *testing.T) {
	assertToken(t, "", emptyRequest())
	assertToken(t, "header_token", tokenInHeader("header_token"))
	assertToken(t, "cookie_token", tokenInCookie("cookie_token"))
}

func tokenInCookie(token string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})
	return req
}

func tokenInHeader(token string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	return req
}

func emptyRequest() *http.Request {
	return httptest.NewRequest(http.MethodPost, "http://example.com/", nil)
}

func assertToken(t *testing.T, token string, req *http.Request) {
	t.Helper()
	assert.Equal(t, token, getToken(req))
}
