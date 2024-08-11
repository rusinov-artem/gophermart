package router

import (
	"net/http"
)

type Handler interface {
	Liveness(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
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

func (r *Router) SetHandler(h Handler) *Router {
	r.RegisterLiveness(h.Liveness)
	r.RegisterRegister(h.Register)
	return r
}
