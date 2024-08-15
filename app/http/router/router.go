package router

import (
	"net/http"
)

type Handler interface {
	Liveness(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	AddOrder(http.ResponseWriter, *http.Request)
	ListOrders(http.ResponseWriter, *http.Request)
	GetBalance(http.ResponseWriter, *http.Request)
	Withdraw(http.ResponseWriter, *http.Request)
}

type Mux interface {
	Method(method, pattern string, handler http.Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Router struct {
	mux Mux
}

func New(mux Mux) *Router {
	r := &Router{
		mux: mux,
	}

	return r
}

func (r *Router) Mux() http.Handler {
	return r.mux
}

func (r *Router) RegisterLiveness(h http.HandlerFunc) {
	r.mux.Method(http.MethodGet, "/liveness", h)
}

func (r *Router) RegisterRegister(h http.HandlerFunc) {
	r.mux.Method(http.MethodPost, "/api/user/register", h)
}

func (r *Router) RegisterLogin(h http.HandlerFunc) {
	r.mux.Method(http.MethodPost, "/api/user/login", h)
}

func (r *Router) RegisterAddOrder(h http.HandlerFunc) {
	r.mux.Method(http.MethodPost, "/api/user/orders", h)
}

func (r *Router) RegisterListOrders(h http.HandlerFunc) {
	r.mux.Method(http.MethodGet, "/api/user/orders", h)
}

func (r *Router) RegisterGetBalance(h http.HandlerFunc) {
	r.mux.Method(http.MethodGet, "/api/user/balance", h)
}

func (r *Router) RegisterWithdraw(h http.HandlerFunc) {
	r.mux.Method(http.MethodPost, "/api/user/balance/withdraw", h)
}

func (r *Router) SetHandler(h Handler) *Router {
	r.RegisterLiveness(h.Liveness)
	r.RegisterRegister(h.Register)
	r.RegisterLogin(h.Login)
	r.RegisterAddOrder(h.AddOrder)
	r.RegisterListOrders(h.ListOrders)
	r.RegisterGetBalance(h.GetBalance)
	r.RegisterWithdraw(h.Withdraw)
	return r
}
