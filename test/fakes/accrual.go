package fakes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/rusinov-artem/gophermart/app/dto"
)

type Accrual struct {
	sync.Mutex
	Server      *httptest.Server
	t           *testing.T
	handlerFunc http.HandlerFunc
	Req         Request
}

type Request struct {
	Method  string
	Path    string
	Headers http.Header
}

func NewAccrual(t *testing.T) *Accrual {
	m := &Accrual{}

	m.Server = httptest.NewServer(m.handlerFunc)

	m.handlerFunc = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - accrual not available"))
	}

	serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		defer m.Unlock()
		m.Req.Method = r.Method
		m.Req.Path = r.URL.Path
		m.Req.Headers = r.Header.Clone()
		if m.handlerFunc != nil {
			m.handlerFunc(w, r)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - accrual service not initialized"))
	})

	server := httptest.NewServer(serverHandler)
	m.Server = server
	m.t = t

	return m
}

func (s *Accrual) URL() string {
	s.Lock()
	defer s.Unlock()
	return s.Server.URL
}

func (s *Accrual) Addr() string {
	s.Lock()
	defer s.Unlock()
	return s.Server.Listener.Addr().String()
}

func (s *Accrual) WillReturn204() {
	s.Lock()
	defer s.Unlock()

	s.handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Accrual) WillReturnOrder(o dto.OrderListItem) {
	s.Lock()
	defer s.Unlock()

	s.handlerFunc = func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `
			{
			  "order": "%s",
			  "status": "%s",
			  "accrual": %f
			}
		`,
			o.OrderNr, o.Status, o.Accrual)
	}
}
